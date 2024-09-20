package main

import (
	"net/http"
	"strconv"

	"github.com/wapo-isai/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't authenticate user")
		return	
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}

	dbChirp, err := cfg.DB.DeleteChirp(chirpID, userID)
	if err != nil {
		respondWithError(w, 403, "Couldn't delete chirp")
		return
	}

	respondWithJSON(w, 204, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
		AuthorId: dbChirp.AuthorId,
	})
}