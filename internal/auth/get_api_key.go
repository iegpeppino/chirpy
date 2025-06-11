package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey -
func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization header")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

// func GetApiKey(header http.Header) (string, error) {

// 	// keyHeader := header.Get("Authorization")
// 	// if keyHeader == "" {
// 	// 	return "", errors.New("couldn't access header")
// 	// }

// 	// sections := strings.SplitN(keyHeader, " ", 2)
// 	// if len(sections) != 2 || sections[0] != "ApiKey" {
// 	// 	return "", errors.New("invalid authorization")
// 	// }
// 	// return sections[1], nil

// 	// authHeader := header.Get("Authorization")
// 	// if authHeader == "" {
// 	// 	return "", errors.New("no authorization header")
// 	// }
// 	// if strings.HasPrefix(authHeader, "ApiKey") {
// 	// 	api_key := strings.TrimSpace(strings.TrimPrefix(authHeader, "ApiKey"))
// 	// 	if api_key != "" || api_key != " " {
// 	// 		return api_key, nil
// 	// 	}
// 	// 	return "", errors.New("api key value not present")
// 	// }
// 	// return "", errors.New("no ApiKey present in authorization header")

// 	authHeader := strings.TrimSpace(header.Get("Authorization"))
// 	if authHeader == "" {
// 		return "", errors.New("no authorization header present")
// 	}
// 	// Normalize multiple spaces to a single space
// 	fields := strings.Fields(authHeader)
// 	if len(fields) != 2 {
// 		return "", errors.New("no key present")
// 	}

// 	// Case-insensitive match for "ApiKey"
// 	if strings.EqualFold(fields[0], "ApiKey") {
// 		if fields[0] != "ApiKey" {
// 			return "", errors.New("field is not ApiKey")
// 		}
// 		return fields[1], nil
// 	}

// 	return "", errors.New("couldn't parse ApiKey")
// }
