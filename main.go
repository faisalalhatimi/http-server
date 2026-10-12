package main

import (
	"log"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}

func main() {
	port := "8080"
	fileServer := http.FileServer(http.Dir("."))

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handlerFunc)
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	server.ListenAndServe()

	log.Printf("Serving files on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
