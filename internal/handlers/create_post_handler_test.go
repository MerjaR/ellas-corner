package handlers

import (
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"ellas-corner/internal/db"
	"ellas-corner/internal/repository"
)

func buildMultipartForm(t *testing.T, fields map[string]string) (*strings.Reader, string) {
	var b strings.Builder
	w := multipart.NewWriter(&b)

	for key, val := range fields {
		err := w.WriteField(key, val)
		if err != nil {
			t.Fatalf("Error writing field %s: %v", key, err)
		}
	}

	err := w.Close()
	if err != nil {
		t.Fatalf("Error closing multipart writer: %v", err)
	}
	return strings.NewReader(b.String()), w.FormDataContentType()
}

func TestCreatePostHandler_POST(t *testing.T) {
	// Step 1: Setup in-memory DB and run migrations
	conn, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}
	if err := conn.RunMigrations(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	repository.SetDatabase(conn)

	// Step 2: Register user
	registerForm := url.Values{}
	registerForm.Set("username", "testuser")
	registerForm.Set("email", "test@example.com")
	registerForm.Set("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(registerForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	RegisterHandler(rr, req)

	// Step 3: Login
	loginForm := url.Values{}
	loginForm.Set("email", "test@example.com")
	loginForm.Set("password", "password123")

	loginReq := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginForm.Encode()))
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginRR := httptest.NewRecorder()
	LoginHandler(loginRR, loginReq)

	var sessionCookie *http.Cookie
	for _, cookie := range loginRR.Result().Cookies() {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("Session token cookie was not set")
	}

	// Step 4: Create a new post
	fields := map[string]string{
		"title":            "Test Integration Post",
		"content":          "This post was created in a test.",
		"category":         "General",
		"is_donation":      "false",
		"donation_country": "",
	}
	body, contentType := buildMultipartForm(t, fields)

	postReq := httptest.NewRequest(http.MethodPost, "/create-post", body)
	postReq.Header.Set("Content-Type", contentType)

	postReq.AddCookie(sessionCookie)
	postRR := httptest.NewRecorder()

	CreatePostHandler(postRR, postReq)

	// Step 5: Check DB for inserted post
	row := conn.Conn.QueryRow("SELECT title, content FROM posts WHERE title = ?", "Test Integration Post")
	var title, content string
	err = row.Scan(&title, &content)
	if err != nil {
		t.Fatalf("Post not found in DB: %v", err)
	}
	if title != "Test Integration Post" || content != "This post was created in a test." {
		t.Errorf("Unexpected values. Got title=%s, content=%s", title, content)
	}
}
