package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain-text password using bcrypt
func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword hashes the password using a cost of 14 (default is 10)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a plain-text password with the hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
