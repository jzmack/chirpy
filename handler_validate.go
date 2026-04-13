package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	type returnVals struct {
		CleanBody string `json:"cleaned_body"`
	}
	respBody := returnVals{
		CleanBody: cleanBody(params.Body),
	}
	respondWithJSON(w, http.StatusOK, respBody)

}

func cleanBody(msg string) string {
	wordList := strings.Split(msg, " ")
	for i := 0; i < len(wordList); i++ {
		lowered := strings.ToLower(wordList[i])
		if lowered == "kerfuffle" || lowered == "sharbert" || lowered == "fornax" {
			wordList[i] = "****"
		}
	}
	return strings.Join(wordList, " ")
}
