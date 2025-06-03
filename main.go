package main

import (
	"log"
	"net/http"
)

func main() {

	const port = "8080"

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./app/"))

	mux.Handle("/app/", http.StripPrefix("/app", fs))
	mux.HandleFunc("/healthz", readinessHandler)

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
	w.Write([]byte("OK"))
}
