package main

type Status string

const (
	Available    Status = "available"
	NotAvailable Status = "not available"
	NotChecking  Status = "not checking"
)

type Link struct {
	URL    string `json:"url"`
	Status Status `json:"status"`
}

type Task struct {
	ID    int    `json:"link_list"`
	Links []Link `json:"links"`
}
