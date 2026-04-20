package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("author_id")
	var allChirps []Chirp
	if idString == "" {
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
	} else {
		parsedID, err := uuid.Parse(idString)
		if err != nil {
			log.Printf("Error parsing author_id: %s", err)
			respondWithError(w, http.StatusBadRequest, "Error parsing author_id")
			return
		}
		dbChirps, err := cfg.db.GetChirpsByAuthor(r.Context(), parsedID)
		if err != nil {
			log.Printf("Error getting chirps from author")
			respondWithError(w, http.StatusInternalServerError, "Error getting chirps from author")
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

}

func (cfg *apiConfig) handlerGetChirpsID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing chirpID")
		return
	}

	dbChirp, err := cfg.db.GetSingleChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirpID not found")
		return
	}

	chirp := Chirp{
		dbChirp.ID,
		dbChirp.CreatedAt,
		dbChirp.UpdatedAt,
		dbChirp.Body,
		dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
