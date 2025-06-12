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
	postIDStr := r.FormValue("id")
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

	categories, err := repository.FetchCategories()
	if err != nil {
		log.Println("Error fetching categories:", err)
		http.Error(w, "Error loading categories", http.StatusInternalServerError)
		return
	}

	// Handle GET request: Show the edit form
	if r.Method == http.MethodGet {
		// Check for user session (if navbar needs it)
		var isLoggedIn bool
		var profilePicture string
		cookie, err := r.Cookie("session_token")
		if err == nil {
			userID, err := repository.GetUserIDBySession(cookie.Value)
			if err == nil && userID != 0 {
				isLoggedIn = true
				user, err := repository.GetUserByID(userID)
				if err == nil {
					profilePicture = user.ProfilePicture
				}
			}
		}

		tmpl, err := template.ParseFiles("web/templates/edit_post.html", "web/templates/partials/navbar.html")
		if err != nil {
			log.Println("Error loading edit template:", err)
			http.Error(w, "Error loading page", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Post":           post,
			"isLoggedIn":     isLoggedIn,
			"ProfilePicture": profilePicture,
			"Categories":     categories,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println("Error executing template:", err)
		}
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
