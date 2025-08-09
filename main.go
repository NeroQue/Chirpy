package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/NEROQUE/Chirpy/handlers"
	"github.com/NEROQUE/Chirpy/internal/database"
	"github.com/NEROQUE/Chirpy/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	filepathRoot := "."
	port := "8080"
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)

	hits := &atomic.Int32{}

	adminCfg := &handlers.AdminConfig{
		FileserverHits: hits,
		DbQueries:      dbQueries,
		Platform:       platform,
		Secret:         os.Getenv("SECRET"),
	}

	metricsMiddleware := middleware.MetricsMiddleware(hits)

	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app", metricsMiddleware(http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlers.Health)
	mux.HandleFunc("GET /api/chirps", adminCfg.HandleGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", adminCfg.HandleGetChirp)
	mux.HandleFunc("POST /api/users", adminCfg.UserHandler)
	mux.HandleFunc("POST /api/chirps", adminCfg.HandleCreateChirps)
	mux.HandleFunc("POST /api/login", adminCfg.HandleLogin)
	mux.HandleFunc("POST /api/refresh", adminCfg.HandleRefreshTokens)
	mux.HandleFunc("POST /api/revoke", adminCfg.HandleRevokeRefreshToken)
	mux.HandleFunc("PUT /api/users", adminCfg.UserUpdateHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", adminCfg.HandleDeleteChirp)

	mux.HandleFunc("GET /admin/metrics", adminCfg.HitHandler)
	mux.HandleFunc("POST /admin/reset", adminCfg.ResetHitsHandler)

	mux.HandleFunc("POST /api/polka/webhooks", adminCfg.PolkaHandler)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Server starting at port %s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server error: %s", err)
	}
}
