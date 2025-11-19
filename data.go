package main

import (
	"encoding/json"
	"fmt"

	"os"
	"sync"
)

var (
	mutex sync.RWMutex

	dataStore []Task
	fileName  = "data.json"
)

// Чтение JSON
func loadData() {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			dataStore = []Task{}
			return
		}
		fmt.Println(err)
		return
	}

	_ = json.Unmarshal(fileData, &dataStore)
}

func saveData() {
	jsonData, _ := json.MarshalIndent(dataStore, "", "  ")
	_ = os.WriteFile(fileName, jsonData, 0644)
}

// Добавление task
func addLink(task Task) Task {
	mutex.Lock()
	defer mutex.Unlock()

	nextID := 1

	if len(dataStore) > 0 {
		lastTask := dataStore[len(dataStore)-1]
		nextID = lastTask.ID + 1
	}

	task.ID = nextID

	dataStore = append(dataStore, task)
	saveData()

	return task
}

// Изменение task
func updateTask(task Task) {
	mutex.Lock()
	defer mutex.Unlock()

	foundIndex := -1
	for i := range dataStore {
		if dataStore[i].ID == task.ID {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return
	}

	dataStore[foundIndex].Links = task.Links

	saveData()
}
