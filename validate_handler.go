package main

import (
	"encoding/json"
	"net/http"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	type requestParams struct {
		Body string `json:"body"`
	}

	// Decode POST request body
	decoder := json.NewDecoder(r.Body)
	params := requestParams{}

	// Call error response helper if error raises during decoding of request
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Validate if Chirp is of 140 or less

	// If not valid call send error helper func
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// If valid, send response and Code 200 with JSON helper func
	sendRespondJSON(w, http.StatusOK, returnVals{Valid: true})
}
