package controller

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"todolist/service"
)

type TodoController struct {
	service *service.TodoServiceImpl
}

func NewTodoController(service *service.TodoServiceImpl) *TodoController {
	return &TodoController{service: service}
}

func (c *TodoController) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	err := c.service.Create(r)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(http.StatusCreated)
}
func (c *TodoController) Get(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	todo, err := c.service.Get(p)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}
func (c *TodoController) Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	err := c.service.Update(r)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(http.StatusOK)
}
func (c *TodoController) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	todos, err := c.service.List(r)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (c *TodoController) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := c.service.Delete(ps)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(http.StatusOK)
}
