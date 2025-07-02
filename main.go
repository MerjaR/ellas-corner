package main

import (
	"log"
	"net/http"

	"ellas-corner/internal/db"
	"ellas-corner/internal/handlers"
	"ellas-corner/internal/repository"
)

func main() {
	// Initialise the SQLite database
	dbInstance, err := db.InitDB("data/forum.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	repository.SetDatabase(dbInstance)
	if err := dbInstance.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create router
	mux := http.NewServeMux()

	// Serve static files from "web/static" when requested at "/static/..." (web added to keep frontend assets in one place)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Set up route handlers

	// Home page and about page
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/about", handlers.AboutHandler)

	// Posts
	mux.HandleFunc("/api/posts", handlers.PostsHandler)
	mux.HandleFunc("/create-post", handlers.CreatePostHandler)
	mux.HandleFunc("/delete-post", handlers.DeletePostHandler)
	mux.HandleFunc("/edit-post", handlers.EditPostHandler)

	// Authentication
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)

	// User
	mux.HandleFunc("/accept-cookies", handlers.AcceptCookiesHandler)
	mux.HandleFunc("/profile", handlers.ProfileHandler)
	mux.HandleFunc("/upload-profile-picture", handlers.UploadProfilePictureHandler)
	mux.HandleFunc("/liked-posts", handlers.LikedPostsHandler)
	mux.HandleFunc("/update-profile-settings", handlers.UpdateProfileSettingsHandler)

	//Comments and reactions
	mux.HandleFunc("/add-comment", handlers.AddCommentHandler)
	mux.HandleFunc("/react", handlers.ReactionHandler)
	mux.HandleFunc("/react-comment", handlers.CommentReactionHandler)
	mux.HandleFunc("/delete-comment", handlers.DeleteCommentHandler)

	//Filtering and search
	mux.HandleFunc("/filter", handlers.FilterHandler)
	mux.HandleFunc("/search", handlers.SearchHandler)

	// Start the server on port 8080
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
