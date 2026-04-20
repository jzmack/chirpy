package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jzmack/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error extracting API key")
		respondWithError(w, http.StatusUnauthorized, "Error extracting API key")
		return
	}
	if apiKey != cfg.polkaKey {
		log.Printf("someone tried to use an unauthorized API Key")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized request")
		return
	}
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters")
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.GetUserByID(r.Context(), params.Data.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error looking up user in DB: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error looking up user in DB")
		return
	}
	err = cfg.db.UpgradeToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		log.Printf("Error upgrading user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error upgrading user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
