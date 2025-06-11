package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSessionToken generates a random token for user sessions
func GenerateSessionToken() string {
	token := make([]byte, 16)
	rand.Read(token)
	return hex.EncodeToString(token)
}
