package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {

	const port = "8080"

	mux := http.NewServeMux()
	apiCfg := apiConfig{}

	fs := http.FileServer(http.Dir("./app/"))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs)))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.getServerRequestsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetServerHits)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := srv.ListenAndServe()
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

func (c *apiConfig) resetServerHits(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	c.fileserverHits.Store(int32(0))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
