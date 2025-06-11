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
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

// Handles User login endpoint
func (c *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {

	// Setting JSON response structure
	type response struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
	}
	decoder := json.NewDecoder(r.Body)

	// Creating frame for request Parameters
	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decoding parameters into frame struct
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Query user with email sent in request
	user, err := c.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get user data", err)
		return
	}

	// Validate user password
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
		// User: User{
		// 	ID:        user.ID,
		// 	CreatedAt: user.CreatedAt,
		// 	UpdatedAt: user.UpdatedAt,
		// 	Email:     user.Email,
		// },
		Email:        user.Email,
		Password:     user.HashedPassword,
		Token:        tokenStr,
		RefreshToken: refreshTokenStr,
		IsChirpyRed:  user.IsChirpyRed,
	})
}

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
