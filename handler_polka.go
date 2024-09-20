package main

import (
	"encoding/json"
	"net/http"

	"github.com/wapo-isai/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolka(w http.ResponseWriter, r *http.Request) {
	type PolkaData struct {
		UserId int 	`json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data PolkaData `json:"data"`
	}

	polkaApiKey, err := auth.GetApiKey(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find API KEY")
		return
	}

	if polkaApiKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find API KEY")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, "success")
		return
	}

	err = cfg.DB.UpdateUserMembership(params.Data.UserId, true)

	if err != nil {
		respondWithError(w, 404, "not found")
		return
	}
	
	respondWithJSON(w, 204, "success")
}