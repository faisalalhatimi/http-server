package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/faisalalhatimi/http-server/internal/auth"
	"github.com/faisalalhatimi/http-server/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) addUserHandler(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	rq := requestBody{}
	if err := decoder.Decode(&rq); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("Database Error: %s", err)
		return
	}

	hashed_password, err := auth.HashPassword(rq.Password)
	if err != nil {
		log.Printf("Hashing Error: %s", err)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          rq.Email,
		HashedPassword: hashed_password,
	})
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

func (cfg *apiConfig) userLoginHandler(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		Expiration *int   `json:"expire_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	rq := requestBody{}
	if err := decoder.Decode(&rq); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("Decoder Error: %s", err)
		return
	}

	defaultExperation := time.Hour
	requestedTime := defaultExperation
	if rq.Expiration != nil {
		requestedTime = time.Duration(*rq.Expiration) * time.Second
	}
	if requestedTime > defaultExperation {
		requestedTime = defaultExperation
	}

	dbUser, err := cfg.db.GetUser(r.Context(), rq.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		log.Printf("Database Error: %s", err)
		return
	}

	match, err := auth.CheckPasswordHash(rq.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("Hashing Error: %s", err)
		return
	}

	if match == false {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		log.Print("Invalid Password")
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, requestedTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("Token couldn't be generated: %s", err)
		return
	}

	response := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Token:     token,
	}
	respondWithJSON(w, http.StatusOK, response)
}
