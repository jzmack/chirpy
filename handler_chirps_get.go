package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	var allChirps []Chirp
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error getting all Chirps: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for i := 0; i < len(dbChirps); i++ {
		chirp := Chirp{
			dbChirps[i].ID,
			dbChirps[i].CreatedAt,
			dbChirps[i].UpdatedAt,
			dbChirps[i].Body,
			dbChirps[i].UserID,
		}
		allChirps = append(allChirps, chirp)
	}
	respondWithJSON(w, http.StatusOK, allChirps)
}
