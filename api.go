package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func postLinks(c *gin.Context) {
	var links struct {
		Links []string `json:"links"`
	}

	if err := c.BindJSON(&links); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	res := make([]Link, 0, len(links.Links))
	for _, str := range links.Links {
		link := Link{
			URL:    str,
			Status: NotChecking,
		}
		res = append(res, link)
	}

	task := addLink(Task{ID: 0, Links: res})
	addTask(task)

	c.JSON(http.StatusOK, task.ID)
}

func getLinks(c *gin.Context) {
	var list struct {
		LinksList []int `json:"links_list"`
	}

	if err := c.BindJSON(&list); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	var result []Task

	mutex.RLock()
	for _, id := range list.LinksList {
		for _, task := range dataStore {
			if task.ID == id {
				taskCopy := Task{
					ID:    task.ID,
					Links: make([]Link, len(task.Links)),
				}

				copy(taskCopy.Links, task.Links)
				result = append(result, taskCopy)
				break
			}
		}
	}
	mutex.RUnlock()

	if len(result) == 0 {
		c.JSON(404, gin.H{})
		return
	}

	pdfBytes, err := generatePDF(result)
	if err != nil {
		c.JSON(500, gin.H{})
		return
	}

	fileName := fmt.Sprintf("report_%d.pdf", time.Now().Unix())

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)

}

func api() {
	router := gin.Default()

	router.GET("/", getLinks)
	router.POST("/", postLinks)

	router.POST("/pause", func(context *gin.Context) {
		pause = true
	})
	router.POST("/resume", func(context *gin.Context) {
		pause = false
	})
	_ = router.Run("localhost:3333")
}
