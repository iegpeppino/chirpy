package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
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
		Email:       updatedUser.Email,
		UpdatedAt:   updatedUser.UpdatedAt,
		IsChirpyRed: updatedUser.IsChirpyRed,
	})

}

// Webhook that makes user a chirpy_red_user
func (c *apiConfig) upgradeUserHandler(w http.ResponseWriter, r *http.Request) {

	type reqParams struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	headerApiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't parse header", err)
		return
	}

	if headerApiKey != c.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "ApiKey not valid", errors.New("api key doesn't match header key"))
		return
	}

	params := reqParams{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		sendRespondJSON(w, 204, nil)
		return
	}

	_, err = c.db.UpgradeUser(r.Context(), params.Data.UserId)
	if err != nil {
		respondWithError(w, 404, "Couldn't update user, invalid data or user does not exist", err)
		return
	}

	sendRespondJSON(w, 204, nil)
}
