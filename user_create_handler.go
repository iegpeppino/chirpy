package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Handles POST user endpoint to create new user in table
func (c *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	// Set frame to map request parameters
	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Decode parameters into frame struct
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	// Hashes password with auth package func
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't hash password", err)
		return
	}
	// Set newUserParameters for the query
	newUserParams := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Email:          params.Email,
		HashedPassword: hashedPass,
	}
	// Query new user creation
	dbUser, err := c.db.CreateUser(r.Context(), newUserParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create new user", err)
		return
	}
	// Send response if query was succesful
	sendRespondJSON(w, 201,
		User{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.UpdatedAt,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
	)
}
