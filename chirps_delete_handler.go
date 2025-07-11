package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"net/http"

	"github.com/google/uuid"
)

func (c *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Text  string `json:"response_text"`
		Chirp Chirp
	}

	// Authentication (validate bearer token)
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't retrieve access token", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenStr, c.secret)
	if err != nil {
		respondWithError(w, 403, "Couldn't validate access token", err)
		return
	}

	// Parse chirp id from URL
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 403, "Couldn't parse Id field", err)
		return
	}

	// Query get chirp by ID
	chirp, err := c.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Couldn't find chirp", err)
		return
	}

	// Check if user is creator of chirp
	if chirp.UserID != userID {
		respondWithError(w, 403, "Unauthorized user", err)
		return
	}

	// Query delete chirp from table
	err = c.db.DeleteChirpById(r.Context(),
		database.DeleteChirpByIdParams{
			ID:     chirpID,
			UserID: userID,
		})
	if err != nil {
		respondWithError(w, 404, "Couldn't get chirp", err)
		return
	}

	// Success response
	sendRespondJSON(w, 204,
		response{
			Text: "Chirp deleted succesfully",
			Chirp: Chirp{
				ID:   chirp.ID,
				Body: chirp.Body,
			}})

}
