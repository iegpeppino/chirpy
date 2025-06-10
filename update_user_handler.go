package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
)

func (c *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {

	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't retrieve user token", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenStr, c.secret)
	if err != nil {
		respondWithError(w, 401, "Token malformed or missing", err)
		return
	}

	decoder := json.NewDecoder(r.Body)

	newParams := reqParams{}
	err = decoder.Decode(&newParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	newHashedPassword, err := auth.HashPassword(newParams.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	updatedUser, err := c.db.UpdateUser(
		r.Context(),
		database.UpdateUserParams{
			Email:          newParams.Email,
			HashedPassword: newHashedPassword,
			ID:             userID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	sendRespondJSON(w, 200, User{
		Email:     updatedUser.Email,
		UpdatedAt: updatedUser.UpdatedAt,
	})

}
