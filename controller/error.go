package controller

import (
	"encoding/json"
	"net/http"
	"todolist/model"
)

func handleError(w http.ResponseWriter, err *model.ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}
