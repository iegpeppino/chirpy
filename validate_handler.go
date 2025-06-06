package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type returnVals struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

type requestParams struct {
	Body string `json:"body"`
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	// Decode POST request body
	deco := json.NewDecoder(r.Body)
	params := requestParams{}
	// Since all response are JSON we set the header type for all
	w.Header().Set("Content-Type", "application/json")
	err := deco.Decode(&params)
	if err != nil {
		// Creating response body for error
		respBody := returnVals{
			Valid: false,
			Error: "something went wrong",
		}

		// Marshalling response into a JSON
		data, err2 := json.Marshal(respBody)
		if err2 != nil {
			log.Printf("Error marshalling JSON: %s", err2)
			w.WriteHeader(500)
			return
		}
		// Respond with status Code and JSON response
		w.WriteHeader(r.Response.StatusCode)
		w.Write(data)
	}

	// Validate if Chirp is of 140 or less
	if len(params.Body) > 140 {
		respBody := returnVals{
			Valid: false,
			Error: "Chirp is too long",
		}
		// Marshall Response with error telling the chirp is too long
		data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		// Write JSON response and statuscode to client
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	// If validation passes return Valid= true JSON response
	respBody := returnVals{
		Valid: true,
	}
	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}
