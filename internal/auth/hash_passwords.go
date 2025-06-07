package auth

import "golang.org/x/crypto/bcrypt"

// Hashes input password
func HashPassword(password string) (string, error) {
	hshPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hshPass), nil
}

// Compare saved hashed password with input password
func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
