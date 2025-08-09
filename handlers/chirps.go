package handlers

import (
	"encoding/json"
	"errors"
	"github.com/NEROQUE/Chirpy/internal/auth"
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
		Body string `json:"body"`
	}

	type response struct {
		Chirp
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}
	if len(params.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}
	params.Body = ProfaneReplace(params.Body)
	chirp, err := cfg.DbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create chirp", err)
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
		RespondWithError(w, http.StatusInternalServerError, "Failed to get all chirps", err)
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
		err := errors.New("no chirpID provided")
		RespondWithError(w, http.StatusBadRequest, "Please provide a chirpID", err)
		return
	}

	result, err := cfg.DbQueries.GetChirp(r.Context(), uuid.MustParse(chirpID))
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Failed to get chirp", err)
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

func (cfg *AdminConfig) HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {

	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		err := errors.New("no chirpID provided")
		RespondWithError(w, http.StatusBadRequest, "Please provide a chirpID", err)
		return
	}
	headers := r.Header
	token, err := auth.GetBearerToken(headers)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
	}
	userID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
	}
	chirp, err := cfg.DbQueries.GetChirp(r.Context(), uuid.MustParse(chirpID))
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Failed to get chirp", err)
		return
	}
	if chirp.UserID != userID {
		RespondWithError(w, http.StatusForbidden, "Forbidden", err)
		return
	}
	_, err = cfg.DbQueries.DeleteChirp(r.Context(), uuid.MustParse(chirpID))
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Failed to delete chirp, not found", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
