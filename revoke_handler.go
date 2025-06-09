package main

import (
	"chirpy/internal/auth"
	"net/http"
)

// Handles revoke refresh token endpoint
func (c *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {

	// Get user's refresh token
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get token", err)
		return
	}

	// Query the revoke
	_, err = c.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
		return
	}
	// Success
	w.WriteHeader(http.StatusNoContent)
}
