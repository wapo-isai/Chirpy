package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
		AuthorId: dbChirp.AuthorId,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	author_id := r.URL.Query().Get("author_id")
	sorting_order := r.URL.Query().Get("sort")

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
			AuthorId: dbChirp.AuthorId,
		})
	}

	if sorting_order == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	}

	if author_id != "" {
		filteredChirps := []Chirp{}
		
		userId, err := strconv.Atoi(author_id)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
			return
		}

		for _, e := range chirps {
			if e.AuthorId == userId {
				filteredChirps = append(filteredChirps, e)
			}
		}

		respondWithJSON(w, http.StatusOK, filteredChirps)
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
