package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/NEROQUE/Chirpy/internal/database"
)

type AdminConfig struct {
	FileserverHits *atomic.Int32
	DbQueries      *database.Queries
	Platform       string
}

func (cfg *AdminConfig) HitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", cfg.FileserverHits.Load())))
}

func (cfg *AdminConfig) ResetHitsHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		RespondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}
	err := cfg.DbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		log.Fatal(err)
		return
	}

	cfg.FileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset and Users deleted!"))
}
