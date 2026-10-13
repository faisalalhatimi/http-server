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

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	s := replaceProfane(params.Body)

	type response struct {
		Body string `json:"cleaned_body"`
	}
	respondWithJSON(w, http.StatusOK, response{Body: s})
}
