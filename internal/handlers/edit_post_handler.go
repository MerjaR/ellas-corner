package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
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
		sessionUser, err := utils.GetSessionUser(r)
		isLoggedIn := err == nil
		profilePicture := ""
		if isLoggedIn {
			profilePicture = sessionUser.ProfilePicture
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
		err := r.ParseMultipartForm(10 << 20) // Max 10MB
		if err != nil {
			log.Println("Error parsing form:", err)
			http.Error(w, "Invalid form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")
		isDonation := r.FormValue("is_donation") != "on"

		// Check if a new image was uploaded
		var imagePath string
		file, header, err := r.FormFile("image")
		if err == nil && header.Size > 0 {
			defer file.Close()
			imagePath, err = repository.SaveImageFile(file, header)
			if err != nil {
				log.Println("Failed to save image:", err)
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}
		} else {
			imagePath = post.Image // keep existing image
		}

		err = repository.UpdatePostWithImage(postID, title, content, category, isDonation, imagePath)
		if err != nil {
			log.Println("Error updating post:", err)
			http.Error(w, "Error updating post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}
