package model

type TaskRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Status      bool   `json:"status"`
}

type ListRequest struct {
	Limit  *int   `json:"limit"`
	Offset *int   `json:"offset"`
	Date   string `json:"date"`
	Status *bool  `json:"status"`
}
