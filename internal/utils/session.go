package utils

import (
	"crypto/rand"
	"ellas-corner/internal/repository"
	"encoding/hex"
	"errors"
	"net/http"
)

// GenerateSessionToken generates a random token for user sessions
func GenerateSessionToken() string {
	token := make([]byte, 16)
	rand.Read(token)
	return hex.EncodeToString(token)
}

type SessionUser struct {
	ID             int
	Username       string
	ProfilePicture string
	Country        string
}

var ErrUnauthenticated = errors.New("user not authenticated")

func GetSessionUser(r *http.Request) (*SessionUser, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, ErrUnauthenticated
	}

	userID, err := repository.GetUserIDBySession(cookie.Value)
	if err != nil || userID == 0 {
		return nil, ErrUnauthenticated
	}

	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err // genuine DB error
	}

	return &SessionUser{
		ID:             user.ID,
		Username:       user.Username,
		ProfilePicture: user.ProfilePicture,
		Country:        user.Country,
	}, nil
}
