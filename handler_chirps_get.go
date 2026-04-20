package main

import (
	"errors"
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func sortOrderFromRequest(r *http.Request) (string, error) {
	orderString := r.URL.Query().Get("sort")
	if orderString == "" {
		return "", nil
	} else if orderString == "asc" {
		return "asc", nil
	} else if orderString == "desc" {
		return "desc", nil
	}
	return "", errors.New("Not a valid sort by expression")
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("author_id")
	sortBy, err := sortOrderFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Not a valid sort by expression")
	}
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
		if sortBy == "desc" {
			sort.Slice(allChirps, func(i, j int) bool {
				return allChirps[i].CreatedAt.After(allChirps[j].CreatedAt)
			})
			respondWithJSON(w, http.StatusOK, allChirps)
		} else {
			sort.Slice(allChirps, func(i, j int) bool {
				return allChirps[i].CreatedAt.Before(allChirps[j].CreatedAt)
			})
			respondWithJSON(w, http.StatusOK, allChirps)
		}
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
		if sortBy == "desc" {
			sort.Slice(allChirps, func(i, j int) bool {
				return allChirps[i].CreatedAt.After(allChirps[j].CreatedAt)
			})
			respondWithJSON(w, http.StatusOK, allChirps)
		} else {
			sort.Slice(allChirps, func(i, j int) bool {
				return allChirps[i].CreatedAt.Before(allChirps[j].CreatedAt)
			})
			respondWithJSON(w, http.StatusOK, allChirps)
		}
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
