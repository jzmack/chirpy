package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jzmack/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	if passwordCheck == true {
		respondWithJSON(w, http.StatusOK, User{
			dbUser.ID,
			dbUser.CreatedAt,
			dbUser.UpdatedAt,
			dbUser.Email,
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
}
