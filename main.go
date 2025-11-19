package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

// Буфер с task которые надо проверить
var tasks = make(chan Task, 1000)

var pause = false

func backgroundWorker() {
	for task := range tasks {
		for pause {
			time.Sleep(time.Second)
			fmt.Println("Ждумс")
		}
		updateTask(VerificationTask(task))
		fmt.Println("Чёто сделал")
	}
}

// VerificationLink Проверка link
func VerificationLink(url string) Status {

	url = strings.TrimSpace(url)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Head(url)
	if err != nil {
		return NotAvailable
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		return Available
	}

	return NotAvailable
}

// VerificationTask Проверка Task
func VerificationTask(task Task) Task {
	var slice []Link
	for _, link := range task.Links {
		status := VerificationLink(link.URL)
		links := Link{link.URL, status}
		slice = append(slice, links)
	}

	task.Links = slice

	return task
}

// StartBufferTasks Нахождение всех ссылок нуждающихся в проверке
func StartBufferTasks() {
	mutex.RLock()
	defer mutex.RUnlock()

	count := 0
	for _, task := range dataStore {
		// Проверяем, есть ли непроверенные ссылки
		hasPending := false
		for _, link := range task.Links {
			if link.Status == NotChecking {
				hasPending = true
				break
			}
		}

		if hasPending {
			go func(task Task) {
				tasks <- task
			}(task)
			count++
		}
	}
}

// addTask Добавляем новую задачу
func addTask(task Task) {
	tasks <- task
}

// Генерация pdf
func generatePDF(tasks []Task) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Tasks Report")
	pdf.Ln(12)

	for _, task := range tasks {
		pdf.SetFont("Arial", "B", 12)
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(0, 10, fmt.Sprintf("Task #%d", task.ID))
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 10)

		for _, link := range task.Links {
			if link.Status == "available" || link.Status == "Available" {
				pdf.SetTextColor(0, 150, 0)
			} else {
				pdf.SetTextColor(200, 0, 0)
			}

			text := fmt.Sprintf(" - %s [%s]", link.URL, link.Status)

			pdf.Cell(0, 6, text)
			pdf.Ln(6)
		}

		pdf.Ln(4)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func main() {
	loadData()
	go backgroundWorker()
	StartBufferTasks()
	api()
}
