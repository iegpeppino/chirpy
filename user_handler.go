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

	type response struct {
		User         User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
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
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Generate JWT for user
	tokenStr, err := auth.MakeJWT(user.ID, c.secret, time.Hour*1)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate access JWT", err)
		return
	}

	// Generate refresh token
	refreshTokenStr, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh token", err)
		return
	}

	_, err = c.db.GenerateRefreshToken(
		r.Context(),
		database.GenerateRefreshTokenParams{
			Token:     refreshTokenStr,
			CreatedAt: time.Now().UTC(), // These attributes could be set in the
			UpdatedAt: time.Now().UTC(), // createTokenQuery itself and not passed down the first time
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh token", err)
		return
	}

	sendRespondJSON(w, 200, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        tokenStr,
		RefreshToken: refreshTokenStr,
	})
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
		return
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

	sendRespondJSON(w, 201,
		User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
	)
}
