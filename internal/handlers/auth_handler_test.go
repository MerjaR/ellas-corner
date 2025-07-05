package handlers

import (
	"ellas-corner/internal/db"
	"ellas-corner/internal/repository"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func setupTestAuthDB(t *testing.T) {
	conn, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize test DB: %v", err)
	}
	err = conn.RunMigrations()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	repository.SetDatabase(conn)
}

func TestRegisterHandler_POST(t *testing.T) {
	setupTestAuthDB(t)

	form := url.Values{}
	form.Add("username", "testuser")
	form.Add("email", "test@example.com")
	form.Add("password", "secretpass")

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	RegisterHandler(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected redirect status, got %d", rr.Code)
	}
	location := rr.Header().Get("Location")
	if location == "" || !strings.Contains(location, "/login") {
		t.Errorf("expected redirect to /login, got %s", location)
	}

	// Check DB entry exists
	user, err := repository.GetUserByEmail("test@example.com")
	if err != nil || user == nil {
		t.Fatalf("expected user to be created in DB, got err=%v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %s", user.Username)
	}
	if user.Password == "secretpass" {
		t.Errorf("expected hashed password, got plain-text")
	}
}
