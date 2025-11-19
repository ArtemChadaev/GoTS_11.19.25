package main

import (
	"net/http"

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

	task := Task{ID: 0, Links: res}
	addLink(task)
	addTask(task)
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
	defer mutex.RUnlock()

	for _, id := range list.LinksList {
		for _, task := range dataStore {
			if task.ID == id {
				result = append(result, task)
				break
			}
		}
	}

	c.JSON(http.StatusOK, result)
}

func api() {
	router := gin.Default()

	router.GET("/", getLinks)
	router.POST("/", postLinks)

	_ = router.Run("localhost:3333")
}
