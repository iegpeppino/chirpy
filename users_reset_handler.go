package main

import "net/http"

// Deletes all users from table "users"
func (c *apiConfig) resetServerUsers(w http.ResponseWriter, r *http.Request) {

	// You must be in "dev" platform in order to be able
	// to run this CRUD operation
	if c.platform != "dev" {
		w.WriteHeader(403)
		w.Write([]byte("403 Forbidden\n"))
		return
	}

	c.fileserverHits.Store(int32(0)) // Resets server hits counter

	// Delete from users query
	err := c.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't delete users from table", err)
		return
	}

	// Success
	w.WriteHeader(200)
}
