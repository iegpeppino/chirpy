package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (c *apiConfig) newChirpHandler(w http.ResponseWriter, r *http.Request) {

	type reqVals struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	// Decode request JSON and map to struct
	decoder := json.NewDecoder(r.Body)
	reqParams := reqVals{}

	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Validate length and proper language
	if len(reqParams.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is longer than 140 characters", nil)
		return
	}

	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	uncleanBody := strings.Split(reqParams.Body, " ")
	for i, word := range uncleanBody {
		for _, bWord := range badWords {
			if strings.ToLower(word) == bWord {
				uncleanBody[i] = "****"
			}
		}
	}

	cleanBody := strings.Join(uncleanBody, " ")

	newChirp := Chirp{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Body:      cleanBody,
		UserID:    reqParams.UserId,
	}

	_, err = c.db.CreateChirp(r.Context(), database.CreateChirpParams(newChirp))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create chirp", err)
		return
	}

	sendRespondJSON(w, 201, newChirp)
}
