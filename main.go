package main

import (
	"io"
	"net/http"
	"time"
)

// Буфер с task которые надо проверить
var tasks = make(chan Task, 100)

func backgroundWorker() {
	for task := range tasks {
		updateTask(VerificationTask(task))
	}
}

// VerificationLink Проверка link
func VerificationLink(url string) Status {
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

func main() {
	loadData()
	go backgroundWorker()
	StartBufferTasks()
	api()
}
