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
}

func main() {

	// Load .env and get env variables to set db connection
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	env_platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("couldn't open connection to database %s", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		db:       dbQueries,
		platform: env_platform,
	}

	const port = "8080"
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./app/"))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs)))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.getServerRequestsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetServerUsers)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.newChirpHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getAllChirpsHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe :", err)
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

func (c *apiConfig) resetServerUsers(w http.ResponseWriter, r *http.Request) {
	if c.platform != "dev" {
		w.WriteHeader(403)
		w.Write([]byte("403 Forbidden\n"))
		return
	}
	w.WriteHeader(200)
	c.fileserverHits.Store(int32(0))
	err := c.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't delete users from table", err)
		return
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
