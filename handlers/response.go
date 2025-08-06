package handlers

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	respBody := map[string]interface{}{
		"error": msg,
	}
	data, _ := json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
