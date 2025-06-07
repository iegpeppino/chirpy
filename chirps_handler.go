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

// Create NewChirp in Db and send JSON response
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

	// Set new chirp parameters
	newChirp := Chirp{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Body:      cleanBody,
		UserID:    reqParams.UserId,
	}

	// Query new chirp creation in db
	_, err = c.db.CreateChirp(r.Context(), database.CreateChirpParams(newChirp))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create chirp", err)
		return
	}
	// If creation is succesful, send JSON response
	sendRespondJSON(w, 201, newChirp)
}

// Lists all chirps in table
func (c *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {

	allChirps, err := c.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	respChirps := mapDbChirpsToJSON(allChirps)

	sendRespondJSON(w, 200, respChirps)

}

// Gets single chirp by ID
func (c *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {

	// Parse id from path to UUID element
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID format", err)
		return
	}

	// Get chirp from db
	chirp, err := c.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Couldn't get chirp", err)
		return
	}

	// Map chirp to response adequate struct
	respChirp := Chirp{
		ID:        chirpID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	sendRespondJSON(w, 200, respChirp)
}

// Maps database.Chirp struct to a JSON ready Chirp struct
func mapDbChirpsToJSON(dbChirps []database.Chirp) []Chirp {
	jsonChirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		jsonChirps[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}
	return jsonChirps
}
