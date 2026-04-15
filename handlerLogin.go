package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jzmack/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Expiry   *int   `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding username and password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding")
		return
	}

	dbUser, err := cfg.db.LoginCheck(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error - query failed: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	passwordCheck, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		log.Printf("Error checking password: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	expiresInSeconds := 3600
	if passwordCheck == true {
		if params.Expiry != nil && *params.Expiry < 3600 {
			expiresInSeconds = *params.Expiry
		}
		duration := time.Duration(expiresInSeconds) * time.Second
		token, err := auth.MakeJWT(dbUser.ID, cfg.apiSecret, duration)
		if err != nil {
			log.Printf("Error generating token: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Error generating token.")
			return
		}
		respondWithJSON(w, http.StatusOK, response{
			User: User{
				dbUser.ID,
				dbUser.CreatedAt,
				dbUser.UpdatedAt,
				dbUser.Email,
			},
			Token: token,
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
}
