package handlers

import (
	"encoding/json"
	"github.com/NEROQUE/Chirpy/internal/database"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

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

func (cfg *AdminConfig) HandleCreateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp
	}

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
	params.Body = ProfaneReplace(params.Body)
	chirp, err := cfg.DbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}
	resp := response{Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}}
	RespondWithJSON(w, http.StatusCreated, resp)
}

func (cfg *AdminConfig) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	var allChirps []Chirp
	result, err := cfg.DbQueries.GetAllChirps(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get all chirps")
		return
	}

	for _, chirp := range result {
		allChirps = append(allChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	RespondWithJSON(w, http.StatusOK, allChirps)
}

func (cfg *AdminConfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		RespondWithError(w, http.StatusBadRequest, "Please provide a chirpID")
		return
	}

	result, err := cfg.DbQueries.GetChirp(r.Context(), uuid.MustParse(chirpID))
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Failed to get chirp")
		return
	}
	chirp := Chirp{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Body:      result.Body,
		UserID:    result.UserID,
	}

	RespondWithJSON(w, http.StatusOK, chirp)
}
