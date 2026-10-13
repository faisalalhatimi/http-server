package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/faisalalhatimi/http-server/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Structs
type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

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

// Handlers
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	w.Write([]byte(body))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "403 Forbidden")
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.db.ResetUsers(r.Context())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data in database has been removed"))
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

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Database Error: %s", err)
	}
	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		db:       dbQueries,
		platform: os.Getenv("PLATFORM"),
	}

	port := "8080"
	fileServer := http.FileServer(http.Dir("."))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /api/healthz", statusHandler)
	mux.HandleFunc("POST /api/validate_chirp", jsonHandler)
	mux.HandleFunc("POST /api/users", apiCfg.addUserHandler)

	// Admin Endpoints
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	server.ListenAndServe()

	log.Fatal(server.ListenAndServe())
}
