package main

import (
	"encoding/json"
	"net/http"
)

func HandlerValidateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	respBody := map[string]interface{}{}
	err := decoder.Decode(&params)
	if err != nil {
		respBody["error"] = "Something went wrong"
		data, _ := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}
	if len(params.Body) > 140 {
		respBody["error"] = "Chirp is too long"
		data, _ := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}
	respBody["valid"] = true
	data, _ := json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}
