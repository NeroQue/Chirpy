package handlers

import (
	"encoding/json"
	"github.com/NEROQUE/Chirpy/internal/auth"
	"github.com/NEROQUE/Chirpy/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *AdminConfig) UserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user, err := cfg.DbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	RespondWithJSON(w, http.StatusCreated, response{User: User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}})
}
