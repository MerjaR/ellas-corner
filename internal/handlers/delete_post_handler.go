package handlers

import (
	"ellas-corner/internal/repository"
	"log"
	"net/http"
	"strconv"
)

// DeletePostHandler handles the deletion of a post
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get the post ID from the form data
	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("Error converting post ID:", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Delete the post from the database
	err = repository.DeletePost(postID)
	if err != nil {
		log.Println("Error deleting post:", err)
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	// Redirect to the profile page after deletion
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
