package main

import (
	"chirpy/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	polkaKey       string
}

func main() {

	// Load .env and get env variables to set db connection
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	secretStr := os.Getenv("SECRET")
	env_platform := os.Getenv("PLATFORM")
	polka_key := os.Getenv("POLKA_KEY")
	if polka_key == "" {
		log.Fatal("POLKA_KEY env variable is not set")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("couldn't open connection to database %s", err)
	}

	dbQueries := database.New(db)

	// Set http client Configuration
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       env_platform,
		secret:         secretStr,
		polkaKey:       polka_key,
	}

	const port = "8080"
	mux := http.NewServeMux()

	// Set url to serve static files
	fs := http.FileServer(http.Dir("./app/"))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs)))

	mux.HandleFunc("GET /api/healthz", readinessHandler)                  // Returns 200 code
	mux.HandleFunc("GET /admin/metrics", apiCfg.getServerRequestsHandler) // Counts server hits

	mux.HandleFunc("POST /admin/reset", apiCfg.resetServerUsers) // Deletes all users from table

	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler) // Create user
	mux.HandleFunc("POST /api/login", apiCfg.loginUserHandler)  // Login user
	mux.HandleFunc("PUT /api/users", apiCfg.updateUserHandler)  // Update email and password

	mux.HandleFunc("POST /api/refresh", apiCfg.refreshTokenHandler) // Refresh access token
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeTokenHandler)   // Revoke refresh token

	mux.HandleFunc("POST /api/chirps", apiCfg.newChirpHandler)                // Create chirp
	mux.HandleFunc("GET /api/chirps", apiCfg.getAllChirpsHandler)             // List all chirps
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpByIDHandler)   // Get specific chirp
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirpHandler) // Delete chirp

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.upgradeUserHandler) // Update user subscription field

	// Set http server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Listen for requests
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe :", err) // Locks the server until interruption
	}

}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK\n"))
}

func (c *apiConfig) getServerRequestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	hits := c.fileserverHits.Load()
	fmt.Fprintf(w, `
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`, hits)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
