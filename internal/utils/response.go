package utils

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Success bool   `json:"success"`
}

func JSON(w http.ResponseWriter, code int, message string, success bool, data any) {
	res := ApiResponse{
		Status:  code,
		Success: success,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

}
