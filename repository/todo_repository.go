package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"todolist/model"
	"todolist/model/db"
)

type TodoRepositoryImpl struct {
	db *sql.DB
}
type TodoRepository interface {
	Create(task model.TaskRequest) error
	ChekId(id int) (bool, error)
	Get(id int) (model.TaskResponse, error)
	Update(task model.TaskRequest) error
	Delete(id int) error
	List(req model.ListRequest) ([]model.TaskResponse, error)
}

func NewTodoRepository(db *sql.DB) *TodoRepositoryImpl {
	return &TodoRepositoryImpl{db: db}
}
func isBlank(str string) bool {
	words := strings.Fields(str)
	if len(words) == 0 {
		return false
	}
	return true
}

// проверка на существование ID
func (r *TodoRepositoryImpl) ChekId(id int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS (SELECT 1 FROM task WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return exists, err
	}
	return exists, nil
}

func (r *TodoRepositoryImpl) Create(task model.TaskRequest) error {
	var description interface{}
	if task.Description == "" || !isBlank(task.Description) {
		description = nil
	} else {
		description = task.Description
	}

	q := "INSERT INTO task (title, description, date, status) VALUES ($1, $2, $3, $4) "
	_, err := r.db.Exec(q, task.Title, description, task.Date, task.Status)
	return err
}
func (r *TodoRepositoryImpl) Get(id int) (model.TaskResponse, error) {
	var task db.TaskDB
	err := r.db.QueryRow("SELECT * FROM task WHERE id = $1", id).Scan(&task.ID, &task.Title, &task.Description, &task.Date, &task.Status)
	if err != nil {
		return model.TaskResponse{}, err
	}
	return task.ToResponse(), nil
}
func (r *TodoRepositoryImpl) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM task WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
func (r *TodoRepositoryImpl) Update(task model.TaskRequest) error {
	var description interface{}
	if task.Description == "" || !isBlank(task.Description) {
		description = nil
	} else {
		description = task.Description
	}
	q := "UPDATE task SET title = $1, description = $2, date = $3, status = $4 WHERE id = $5"
	//можно реализовать иначе, зависит от деталей
	_, err := r.db.Exec(q, task.Title, description, task.Date, task.Status, task.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
func (r *TodoRepositoryImpl) List(req model.ListRequest) ([]model.TaskResponse, error) {
	q := "SELECT * FROM task"
	params := []interface{}{}
	conditions := []string{}
	paramCounter := 1

	if req.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", paramCounter))
		params = append(params, *req.Status)
		paramCounter++
	}

	if req.Date != "" {
		conditions = append(conditions, fmt.Sprintf("date = $%d", paramCounter))
		params = append(params, req.Date)
		paramCounter++
	}

	if len(conditions) > 0 {
		q += " WHERE " + strings.Join(conditions, " AND ")
	}

	q += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCounter, paramCounter+1)
	params = append(params, req.Limit, req.Offset)

	rows, err := r.db.Query(q, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var respTask []model.TaskResponse
	for rows.Next() {
		var task db.TaskDB
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Date, &task.Status)
		if err != nil {
			return nil, err
		}

		respTask = append(respTask, task.ToResponse())

	}
	return respTask, nil
}
