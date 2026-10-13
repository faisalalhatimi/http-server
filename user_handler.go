package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) addUserHandler(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	rq := requestBody{}
	if err := decoder.Decode(&rq); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("Database Error: %s", err)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), rq.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User couldn't be added")
		log.Printf("Database Error: %s", err)
		return
	}

	responseUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, responseUser)
}
