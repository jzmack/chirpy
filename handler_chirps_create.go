package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jzmack/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	cleaned_body := cleanBody(params.Body)
	user_id, err := uuid.Parse(params.UserID)
	if err != nil {
		log.Printf("Error parsing user ID: %s", err)
		w.WriteHeader(500)
		return
	}

	chirp_params := database.CreateChirpParams{
		Body:   cleaned_body,
		UserID: user_id,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), chirp_params)
	if err != nil {
		log.Printf("Error creating Chirp: %s", err)
		return
	}

	createdChirp := Chirp{
		chirp.ID,
		chirp.CreatedAt,
		chirp.UpdatedAt,
		chirp.Body,
		chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, createdChirp)

}

func cleanBody(msg string) string {
	wordList := strings.Split(msg, " ")
	for i := 0; i < len(wordList); i++ {
		lowered := strings.ToLower(wordList[i])
		if lowered == "kerfuffle" || lowered == "sharbert" || lowered == "fornax" {
			wordList[i] = "****"
		}
	}
	return strings.Join(wordList, " ")
}
