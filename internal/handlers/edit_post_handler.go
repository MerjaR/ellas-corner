package handlers

import (
	"ellas-corner/internal/repository"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// EditPostHandler shows the edit form for a post
func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get post ID from the query parameters
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("Error converting post ID:", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Fetch the post by its ID
	post, err := repository.GetPostByID(postIDStr)
	if err != nil {
		log.Println("Error fetching post for edit:", err)
		http.Error(w, "Error fetching post", http.StatusInternalServerError)
		return
	}

	// Handle GET request: Show the edit form
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/edit_post.html")
		if err != nil {
			log.Println("Error loading edit template:", err)
			http.Error(w, "Error loading page", http.StatusInternalServerError)
			return
		}

		// Pass the post data to the template
		tmpl.Execute(w, post)
		return
	}

	// Handle POST request: Update the post
	if r.Method == http.MethodPost {
		// Get form values
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")

		// Update the post in the database
		err = repository.UpdatePost(postID, title, content, category)
		if err != nil {
			log.Println("Error updating post:", err)
			http.Error(w, "Error updating post", http.StatusInternalServerError)
			return
		}

		// Redirect to the profile page after updating
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}
