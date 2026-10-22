package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type WebhookData struct {
	UserID uuid.UUID `json:"user_id"`
}

type PolkaWebhook struct {
	Event string      `json:"event"`
	Data  WebhookData `json:"data"`
}

func (cfg *apiConfig) webhookHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	webhook := PolkaWebhook{}
	if err := decoder.Decode(&webhook); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("Decoder Error: %s", err)
		return
	}

	if webhook.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err := cfg.db.UpgreadUser(r.Context(), webhook.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not Found")
		log.Printf("Database Error: %s", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
