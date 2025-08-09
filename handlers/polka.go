package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

func (cfg *AdminConfig) PolkaHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}
	if params.Event != "user.upgraded" {
		RespondWithError(w, http.StatusNoContent, "Invalid request payload", err)
		return
	}
	_, err = cfg.DbQueries.UpgradeUser(r.Context(), params.Data.UserID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Failed to upgrade user, not found", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
