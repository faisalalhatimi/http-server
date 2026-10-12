package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
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

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	body := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset hits counter to 0"))
}

func main() {
	apiCfg := &apiConfig{}
	port := "8080"
	fileServer := http.FileServer(http.Dir("."))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("/healthz", statusHandler)
	mux.HandleFunc("/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/reset", apiCfg.resetHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	server.ListenAndServe()

	log.Printf("Serving files on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
