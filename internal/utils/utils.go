package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendError(w http.ResponseWriter, code int, errorCode int, message string) {
	w.WriteHeader(code)
	errorResponse := ErrorResponse{
		Code:    errorCode,
		Message: message,
	}
	jsonResponse, _ := json.Marshal(errorResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
