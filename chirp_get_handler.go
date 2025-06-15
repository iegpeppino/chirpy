package main

import (
	"net/http"

	"github.com/google/uuid"
)

// Gets single chirp by ID
func (c *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {

	// Parse id from path to UUID element
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID format", err)
		return
	}

	// Get chirp from db
	chirp, err := c.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Couldn't get chirp", err)
		return
	}

	// Map chirp to response adequate struct
	respChirp := Chirp{
		ID:        chirpID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	// Success response
	sendRespondJSON(w, 200, respChirp)
}
