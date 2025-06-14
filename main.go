package main

import (
	"log"
	"net/http"

	"ellas-corner/internal/db"
	"ellas-corner/internal/handlers"
)

func main() {
	// Initialize the SQLite database
	db.InitDB("forum.db")
	db.RunMigrations()

	// Create a new multiplexer (router)
	mux := http.NewServeMux()

	// Handle static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Home page
	mux.HandleFunc("/", handlers.HomeHandler)

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
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
