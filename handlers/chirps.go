package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ProfaneReplace(s string) string {
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(s, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := profaneWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	var respBody = map[string]interface{}{}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if len(params.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	respBody["cleaned_body"] = ProfaneReplace(params.Body)
	RespondWithJSON(w, http.StatusOK, respBody)
}
