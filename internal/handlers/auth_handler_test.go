package handlers

import (
	"ellas-corner/internal/db"
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
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

func TestLoginHandler_POST(t *testing.T) {
	setupTestAuthDB(t)

	// First, create a user
	username := "loginuser"
	email := "login@example.com"
	password := "secret123"
	hashed, _ := utils.HashPassword(password)
	err := repository.CreateUser(username, email, hashed, "1.png")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Simulate login request
	form := strings.NewReader("email=login@example.com&password=secret123")
	req := httptest.NewRequest(http.MethodPost, "/login", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	LoginHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("Expected status %d, got %d", http.StatusSeeOther, resp.StatusCode)
	}

	cookies := resp.Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "session_token" {
			sessionCookie = c
		}
	}
	if sessionCookie == nil {
		t.Error("Expected session_token cookie to be set")
	}
}
