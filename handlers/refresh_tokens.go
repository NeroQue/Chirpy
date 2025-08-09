package handlers

import (
	"github.com/NEROQUE/Chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *AdminConfig) HandleRefreshTokens(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to get refresh token", err)
		return
	}
	user, err := cfg.DbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Failed to get user from refresh token", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.Secret, time.Hour)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create token", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, response{Token: token})
}

func (cfg *AdminConfig) HandleRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to get refresh token", err)
		return
	}
	_, err = cfg.DbQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to revoke refresh token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
