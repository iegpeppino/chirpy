package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (c *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := c.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get user data", err)
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
	}

	sendRespondJSON(w, 200, mapDbUserToResponseUser(user))
}

func (c *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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
	}

	newUserParams := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Email:          params.Email,
		HashedPassword: hashedPass,
	}

	dbUser, err := c.db.CreateUser(r.Context(), newUserParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create new user", err)
		return
	}

	sendRespondJSON(w, 201, mapDbUserToResponseUser(dbUser))
}

// This helper doesn't return the hashed password as a response!
func mapDbUserToResponseUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}
