package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

// Lists all chirps in table
func (c *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {

	// Parse sort and author_id params from URL
	author := r.URL.Query().Get("author_id")
	sorted := r.URL.Query().Get("sort")
	// If there's author_id get all chirps from that author
	if author != "" {
		authorID, err := uuid.Parse(author)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse author_id", err)
			return
		}
		authorChirps, err := c.db.GetChirpsByUserID(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't get chirps", err)
			return
		}

		if sorted == "desc" {

			sort.Slice(authorChirps, func(i, j int) bool {
				return authorChirps[i].CreatedAt.After(authorChirps[j].CreatedAt)
			})

			respChirps := mapDbChirpsToJSON(authorChirps)
			sendRespondJSON(w, 200, respChirps)
		}
		if sorted == "asc" || sorted == "" {
			respChirps := mapDbChirpsToJSON(authorChirps)
			sendRespondJSON(w, 200, respChirps)
		}
		return
	}

	// If there was no author_id list all chirps from database
	allChirps, err := c.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	if sorted == "desc" {
		sort.Slice(allChirps, func(i, j int) bool {
			return allChirps[i].CreatedAt.After(allChirps[j].CreatedAt)
		})
		respChirps := mapDbChirpsToJSON(allChirps)
		sendRespondJSON(w, 200, respChirps)
	}

	if sorted == "asc" || sorted == "" {
		respChirps := mapDbChirpsToJSON(allChirps)
		sendRespondJSON(w, 200, respChirps)
	}

}
