package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func MakeRefreshToken() (string, error) {
	// Making random bytes
	randBytes := make([]byte, 32)
	// Generating random data
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", errors.New("failed to generate random bytes")
	}

	return hex.EncodeToString(randBytes), nil
}
