package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/wapo-isai/chirpy/internal/auth"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
	AuthorId int `json:"author_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 1. ensure user is authenticated by validating token
	// 2. get the user information from the token
	// 3. save the author_id in the chirp

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't authenticate user")
		return	
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
		AuthorId: chirp.AuthorId,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
