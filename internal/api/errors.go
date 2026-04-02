package api

import (
	"encoding/json"
	"net/http"
)

// APIError is the standard error response shape.
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIError{Status: status, Message: msg})
}

func badRequest(w http.ResponseWriter, msg string) {
	writeError(w, http.StatusBadRequest, msg)
}

func notFound(w http.ResponseWriter, msg string) {
	writeError(w, http.StatusNotFound, msg)
}

func serverError(w http.ResponseWriter, msg string) {
	writeError(w, http.StatusInternalServerError, msg)
}
