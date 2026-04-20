package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jzmack/chirpy/internal/auth"
	"github.com/jzmack/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting access token")
		respondWithError(w, http.StatusUnauthorized, "Error with access token")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.apiSecret)
	if err != nil {
		log.Printf("Error validating token")
		respondWithError(w, http.StatusUnauthorized, "Error validating token")
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	userParams := database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPass,
	}
	user, err := cfg.db.UpdateUser(r.Context(), userParams)

	responseUser := User{
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Email,
		user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, responseUser)
}
