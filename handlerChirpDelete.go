package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jzmack/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error parsing the chirp id")
		respondWithError(w, http.StatusBadRequest, "Error parsing chirp ID")
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting access token from header")
		respondWithError(w, http.StatusUnauthorized, "Error extracting access token")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.apiSecret)
	if err != nil {
		log.Printf("Error validating token")
		respondWithError(w, http.StatusUnauthorized, "Error validating token")
		return
	}
	dbChirp, err := cfg.db.GetSingleChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Could not find chirp, i guess")
		respondWithError(w, http.StatusNotFound, "Could not find chirp")
		return
	}

	if userID != dbChirp.UserID {
		log.Printf("Not authorized - user IDs do not match")
		respondWithError(w, http.StatusForbidden, "Not authorized")
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), dbChirp.ID)
	if err != nil {
		log.Printf("Error deleting chirp")
		respondWithError(w, http.StatusInternalServerError, "Cannot find chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
