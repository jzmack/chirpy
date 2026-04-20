package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jzmack/chirpy/internal/auth"
	"github.com/jzmack/chirpy/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %s", err)
		w.WriteHeader(400)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	user_params := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	}

	user, err := cfg.db.CreateUser(r.Context(), user_params)
	if err != nil {
		log.Printf("Error adding user to database: %s ", err)
	}

	createdUser := User{
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Email,
		user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusCreated, createdUser)
}
