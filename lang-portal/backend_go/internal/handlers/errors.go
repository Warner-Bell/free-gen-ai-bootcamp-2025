package handlers

import (
    "encoding/json"
    "net/http"
)

type ErrorResponse struct {
    Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) error {
    return writeJSON(w, status, ErrorResponse{Error: message})
}
