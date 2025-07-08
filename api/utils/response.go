package utils

import (
	"encoding/json"
	"net/http"
)

func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := map[string]string{
		"error":  message,
		"status": "error",
	}
	SendJSONResponse(w, statusCode, errorResponse)
}