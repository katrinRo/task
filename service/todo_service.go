package service

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
	"time"
	"todolist/model"
	"todolist/repository"
)

type TodoService interface {
	Create(r *http.Request) *model.ErrorResponse
	Get(p httprouter.Params) (model.TaskResponse, *model.ErrorResponse)
	Update(r *http.Request, p httprouter.Params) *model.ErrorResponse
	Delete(r *http.Request) *model.ErrorResponse
	List(r *http.Request) ([]model.TaskResponse, *model.ErrorResponse)
}

type TodoServiceImpl struct {
	repository repository.TodoRepository
}

func NewTodoService(repository repository.TodoRepository) *TodoServiceImpl {
	return &TodoServiceImpl{repository: repository}
}

func isValidCreate(task model.TaskRequest) bool {
	return validTask(task)
}

func isValidId(idStr string) (int, error) {
	return strconv.Atoi(idStr)
}
func chekDate(data string) bool {
	dateLayout := "2006-01-02"
	_, err := time.Parse(dateLayout, data)
	if err != nil {
		return false
	}
	return true
}
func isValidList(req model.ListRequest) bool {
	if req.Limit == nil {
		return false
	}
	if req.Offset == nil {
		return false
	}
	if req.Date != "" {
		if !chekDate(req.Date) {
			return false
		}
	}
	return true
}
func validTask(task model.TaskRequest) bool {
	if task.Title == "" || len(task.Title) > 100 {
		return false
	}
	if len(task.Description) > 400 {
		return false
	}
	if task.Date == "" {
		return false
	}

	return chekDate(task.Date)
}

func isValidUpdate(task model.TaskRequest) bool {
	if task.ID == 0 {
		return false
	}
	return validTask(task)
}

func (s *TodoServiceImpl) Create(r *http.Request) *model.ErrorResponse {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &model.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Не удалось прочитать тело запроса",
		}
	}

	task := model.TaskRequest{}
	err = json.Unmarshal(body, &task)
	if err != nil {
		return &model.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Некорректные данные",
		}
	}
	if !isValidCreate(task) {
		return &model.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Некорректные данные, ошибка валидации",
		}
	}
	err = s.repository.Create(task)
	if err != nil {
		return &model.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Ошибка при создании задачи",
		}
	}
	return nil
}
func (s *TodoServiceImpl) Get(ps httprouter.Params) (model.TaskResponse, *model.ErrorResponse) {
	idStr := ps.ByName("id")
	if idStr == "" {
		return model.TaskResponse{}, &model.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Отсутсвует ID",
		}
	}
	id, err := isValidId(idStr)
	if err != nil {
		return model.TaskResponse{}, &model.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Некорректные данные, ошибка валидации",
		}
	}
	exists, err := s.repository.ChekId(id)
	if err != nil {
		return model.TaskResponse{}, &model.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Ошибка при выполнении запроса",
		}
	}
	if !exists {
		return model.TaskResponse{}, &model.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "ID не существует",
		}
	}
	task, err := s.repository.Get(id)
	if err != nil {
		return task, &model.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Ошибка при выполнении запроса",
		}
	}
	return task, nil
}
func (s *TodoServiceImpl) Update(r *http.Request) *model.ErrorResponse {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &model.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Не удалось прочитать тело запроса",
		}
	}

	task := model.TaskRequest{}
	err = json.Unmarshal(body, &task)
	if err != nil {
		return &model.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Некорректные данные",
		}
	}
	if !isValidUpdate(task) {
		return &model.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Некорректные данные, ошибка валидации",
		}
	}
	exists, err := s.repository.ChekId(task.ID)
	if err != nil {
		return &model.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Ошибка при выполнении запроса",
		}
	}
	if !exists {
		return &model.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "ID не существует",
		}
	}
	err = s.repository.Update(task)
	if err != nil {
		return &model.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Ошибка при выполнении запроса",
		}
	}
	return nil
}
func (s *TodoServiceImpl) Delete(ps httprouter.Params) *model.ErrorResponse {
	idStr := ps.ByName("id")
	if idStr == "" {
		return &model.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Отсутсвует ID",
		}
	}
	id, err := isValidId(idStr)
	if err != nil {
		return &model.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Некорректные данные, ошибка валидации",
		}
	}
	exists, err := s.repository.ChekId(id)
	if err != nil {
		return &model.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Ошибка при выполнении запроса",
		}
	}
	if !exists {
		return &model.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "ID не существует",
		}
	}
	err = s.repository.Delete(id)
	if err != nil {
		return &model.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Ошибка при выполнении запроса",
		}
	}
	return nil
}
func (s *TodoServiceImpl) List(r *http.Request) ([]model.TaskResponse, *model.ErrorResponse) {
	var body model.ListRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, &model.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Не удалось прочитать тело запроса",
		}
	}

	if !isValidList(body) {
		return nil, &model.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Не валидный запрос",
		}
	}
	task, err := s.repository.List(body)
	if err != nil {
		return nil, &model.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Ошибка при выполнении запроса",
		}
	}
	return task, nil
}
