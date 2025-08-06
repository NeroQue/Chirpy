package handlers

import (
	"encoding/json"
	"github.com/NEROQUE/Chirpy/internal/auth"
	"net/http"
)

func (cfg *AdminConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
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

	user, err := cfg.DbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	RespondWithJSON(w, http.StatusOK, response{User: User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}})
}
