package utils

import (
	"encoding/json"
	"net/http"
)

type appResponse struct {
	StatusCode int    `json:"statusCode"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	Data       any    `json:"data,omitempty"`
}

func SendSuccessResponse(w http.ResponseWriter, data any, statusCode int) {
	resp := &appResponse{
		StatusCode: statusCode,
		Success:    true,
		Data:       data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "error encoding success response json", http.StatusInternalServerError)
	}
}

func SendErrorResponse(w http.ResponseWriter, err string, statusCode int) {
	resp := &appResponse{
		StatusCode: statusCode,
		Success:    false,
		Error:      err,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "error encoding error response json", http.StatusInternalServerError)
	}
}
