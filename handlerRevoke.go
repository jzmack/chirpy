package main

import (
	"log"
	"net/http"

	"github.com/jzmack/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token from headers")
		respondWithError(w, http.StatusUnauthorized, "Invalid token in headers")
		return
	}
	dbRefreshToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Error looking up refresh token in DB")
		respondWithError(w, http.StatusUnauthorized, "Refresh token not found")
		return
	}
	err = cfg.db.RevokeToken(r.Context(), dbRefreshToken.Token)
	if err != nil {
		log.Printf("Error revoking token")
		respondWithError(w, http.StatusUnauthorized, "Error revoking token")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
