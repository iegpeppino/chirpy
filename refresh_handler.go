package main

import (
	"chirpy/internal/auth"
	"fmt"
	"net/http"
	"time"
)

// If the user has a valid token, a new refreshed access token is created
func (c *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	fmt.Println(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get token", err)
		return
	}

	user, err := c.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, c.secret, time.Hour*1)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	sendRespondJSON(w, http.StatusOK, response{Token: accessToken})
}
