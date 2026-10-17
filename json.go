package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Helpers
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{Error: msg})
}

// Functions
func replaceProfane(original string) string {
	profane := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(original, " ")

	for i, word := range words {
		lowerd := strings.ToLower(word)
		if _, ok := profane[lowerd]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
