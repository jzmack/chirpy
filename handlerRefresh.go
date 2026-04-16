package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jzmack/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error extracting refresh token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Error extracting refresh token")
		return
	}
	dbRefreshToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Could not find refresh token in database: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Could not find refresh token in database")
		return
	}
	if dbRefreshToken.RevokedAt.Valid {
		log.Printf("User %s tried to refresh with a revoked refresh token", dbRefreshToken.UserID)
		respondWithError(w, http.StatusUnauthorized, "Cannot use revoked token")
		return
	}
	if dbRefreshToken.ExpiresAt.Before(time.Now()) {
		log.Printf("User %s tried to refresh with an expired token", dbRefreshToken.UserID)
		respondWithError(w, http.StatusUnauthorized, "Cannot use expired token")
		return
	}
	tokenExpires := time.Duration(3600) * time.Second
	newToken, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.apiSecret, tokenExpires)
	if err != nil {
		log.Printf("Error generating new token %s", err)
		respondWithError(w, http.StatusUnauthorized, "Error generating new token")
		return
	}
	respondWithJSON(w,
		http.StatusOK,
		response{Token: newToken})
}
