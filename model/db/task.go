package db

import (
	"database/sql"
	"time"
	"todolist/model"
)

type TaskDB struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Date        time.Time      `json:"date"`
	Status      bool           `json:"status"`
}

func (t *TaskDB) StatusString() string {
	if t.Status {
		return "Выполнено"
	}
	return "Не выполнено"
}
func (t *TaskDB) DateFormat() string {
	return t.Date.Format("2006-01-02")
}
func (t *TaskDB) ToResponse() model.TaskResponse {
	return model.TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description.String,
		Date:        t.DateFormat(),
		Status:      t.StatusString(),
	}
}
