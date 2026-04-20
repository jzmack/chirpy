package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jzmack/chirpy/internal/auth"
	"github.com/jzmack/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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
	refreshLength := time.Hour * 60 * 24
	if passwordCheck == true {
		duration := time.Duration(expiresInSeconds) * time.Second
		token, err := auth.MakeJWT(dbUser.ID, cfg.apiSecret, duration)
		if err != nil {
			log.Printf("Error generating access token: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Error generating access token.")
			return
		}

		refreshExpires := time.Now().Add(refreshLength)
		refreshToken := auth.MakeRefreshToken()
		refreshParams := database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    dbUser.ID,
			ExpiresAt: refreshExpires,
			RevokedAt: sql.NullTime{},
		}
		dbRefresh, err := cfg.db.CreateRefreshToken(r.Context(), refreshParams)
		if err != nil {
			log.Printf("Error creating refresh token: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
			return
		}

		respondWithJSON(w, http.StatusOK, response{
			User: User{
				dbUser.ID,
				dbUser.CreatedAt,
				dbUser.UpdatedAt,
				dbUser.Email,
				dbUser.IsChirpyRed,
			},
			Token:        token,
			RefreshToken: dbRefresh.Token,
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
}
