package main

import (
	"chirpy/internal/auth"
	"net/http"
	"time"
)

// If the user has a valid token, a new refreshed access token is created
func (c *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Set response frame struct
	type response struct {
		Token string `json:"token"`
	}

	// Get user's refresh token
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't get token", err)
		return
	}
	// Query user data using the refresh token
	user, err := c.db.GetUserFromRefreshToken(r.Context(), tokenStr)
	if err != nil {
		respondWithError(w, 401, "Couldn't get user for refresh token", err)
		return
	}

	// Create new (refresh) access token
	accessToken, err := auth.MakeJWT(user.ID, c.secret, time.Hour*1)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}
	// Send response for success
	sendRespondJSON(w, http.StatusOK, response{Token: accessToken})
}
