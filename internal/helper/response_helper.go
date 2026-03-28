package helper

import (
	"encoding/json"
	"net/http"
)

func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func SendError(w http.ResponseWriter, message string, status int) {
	SendJSON(w, status, map[string]string{
		"error": message,
	})
}
