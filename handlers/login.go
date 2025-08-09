package handlers

import (
	"encoding/json"
	"github.com/NEROQUE/Chirpy/internal/auth"
	"github.com/NEROQUE/Chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *AdminConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User         User   `json:"user,omitempty"`
		Token        string `json:"token,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
		IsChirpyRed  bool   `json:"is_chirpy_red,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	user, err := cfg.DbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	tokenSecret := cfg.Secret
	expiresIn := time.Hour

	token, err := auth.MakeJWT(user.ID, tokenSecret, expiresIn)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create refresh token", err)
		return
	}

	_, err = cfg.DbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create refresh token in the database", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, response{User: User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	},
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed})
}
