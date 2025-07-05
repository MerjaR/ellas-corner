package utils_test

import (
	"testing"

	"ellas-corner/internal/utils"
)

func TestPasswordHashing(t *testing.T) {
	originalPassword := "MyS3cureP@ssword"
	wrongPassword := "NotMyPassword"

	// Hash the original password
	hash, err := utils.HashPassword(originalPassword)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Ensure the hash is not empty and not equal to the plain password
	if hash == "" {
		t.Error("HashPassword returned an empty hash")
	}
	if hash == originalPassword {
		t.Error("Hashed password should not match the original password")
	}

	// Ensure CheckPasswordHash works for correct password
	if !utils.CheckPasswordHash(originalPassword, hash) {
		t.Error("CheckPasswordHash failed to validate correct password")
	}

	// Ensure CheckPasswordHash fails for incorrect password
	if utils.CheckPasswordHash(wrongPassword, hash) {
		t.Error("CheckPasswordHash incorrectly validated wrong password")
	}
}
