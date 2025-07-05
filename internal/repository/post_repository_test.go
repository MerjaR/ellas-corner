package repository_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"ellas-corner/internal/db"
	"ellas-corner/internal/repository"
)

// Sets up a test database for the test run only, does not affect forum.db
func setupTestDB(t *testing.T) *sql.DB {
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}

	schema := `
CREATE TABLE posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER,
	title TEXT,
	content TEXT,
	category TEXT,
	image TEXT,
	is_donation BOOLEAN,
	donation_country TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT,
	profile_picture TEXT
);
CREATE TABLE post_reactions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER,
	user_id INTEGER,
	reaction_type TEXT
);`

	_, err = conn.Exec(schema)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	// Inject the test DB into the repository
	repository.SetDatabase(&db.Database{Conn: conn})

	return conn
}

func TestCreatePost(t *testing.T) {
	testConn := setupTestDB(t)
	defer testConn.Close()

	userID := 1
	_, err := testConn.Exec("INSERT INTO users (id, username, profile_picture) VALUES (?, ?, ?)", userID, "ella", "pic.jpg")
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	err = repository.CreatePost(userID, "Test Title", "Test Content", "General", "image.jpg", true, "Germany")
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	row := testConn.QueryRow("SELECT title FROM posts WHERE user_id = ?", userID)
	var title string
	err = row.Scan(&title)
	if err != nil {
		t.Fatalf("failed to read inserted post: %v", err)
	}
	if title != "Test Title" {
		t.Errorf("expected 'Test Title', got '%s'", title)
	}
}

func TestAddReaction(t *testing.T) {
	testConn := setupTestDB(t)
	defer testConn.Close()

	// Create test user and post
	userID := 1
	postID := 1

	_, err := testConn.Exec("INSERT INTO users (id, username, profile_picture) VALUES (?, ?, ?)", userID, "ella", "pic.jpg")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}
	_, err = testConn.Exec("INSERT INTO posts (id, user_id, title, content, category) VALUES (?, ?, ?, ?, ?)", postID, userID, "Test Post", "Some content", "General")
	if err != nil {
		t.Fatalf("failed to insert post: %v", err)
	}

	// 1️⃣ Insert new "like" reaction
	err = repository.AddReaction(userID, postID, "like")
	if err != nil {
		t.Fatalf("AddReaction failed on insert: %v", err)
	}

	var reaction string
	err = testConn.QueryRow("SELECT reaction_type FROM post_reactions WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&reaction)
	if err != nil {
		t.Fatalf("failed to fetch reaction: %v", err)
	}
	if reaction != "like" {
		t.Errorf("expected 'like', got '%s'", reaction)
	}

	// 2️⃣ Change to "dislike"
	err = repository.AddReaction(userID, postID, "dislike")
	if err != nil {
		t.Fatalf("AddReaction failed on update: %v", err)
	}

	err = testConn.QueryRow("SELECT reaction_type FROM post_reactions WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&reaction)
	if err != nil {
		t.Fatalf("failed to fetch updated reaction: %v", err)
	}
	if reaction != "dislike" {
		t.Errorf("expected 'dislike', got '%s'", reaction)
	}

	// 3️⃣ React again with "dislike" (same) → should be no-op but still same value
	err = repository.AddReaction(userID, postID, "dislike")
	if err != nil {
		t.Fatalf("AddReaction failed on no-op: %v", err)
	}

	err = testConn.QueryRow("SELECT reaction_type FROM post_reactions WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&reaction)
	if err != nil {
		t.Fatalf("failed to fetch reaction after no-op: %v", err)
	}
	if reaction != "dislike" {
		t.Errorf("expected 'dislike' after no-op, got '%s'", reaction)
	}
}
