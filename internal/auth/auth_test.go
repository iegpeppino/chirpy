package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// func TestCheckPasswordHash(t *testing.T) {
// 	pass1 := "superpassword123!"
// 	pass2 := "megapassword543!"
// 	hash1, _ := HashPassword(pass1)
// 	hash2, _ := HashPassword(pass2)

// 	tests := []struct {
// 		name     string
// 		password string
// 		hash     string
// 		wantErr  bool
// 	}{
// 		{
// 			name:     "Correct Password",
// 			password: pass1,
// 			hash:     hash1,
// 			wantErr:  false,
// 		},
// 		{
// 			name:     "Incorrect Password",
// 			password: "WhateverPassword",
// 			hash:     hash1,
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "Password doesn't match different hash",
// 			password: pass1,
// 			hash:     hash2,
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "Empty Password",
// 			password: "",
// 			hash:     hash1,
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "Invalid Hash",
// 			password: pass1,
// 			hash:     "thishashisnotcorrect",
// 			wantErr:  true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := CheckPasswordHash(tt.password, tt.hash)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour*1)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid Token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid Token",
			tokenString: "mamma.mia.whatisdis",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong Secret",
			tokenString: validToken,
			tokenSecret: "knownfact",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error= %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
