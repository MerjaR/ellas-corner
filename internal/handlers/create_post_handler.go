package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in (session token)
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// Redirect to login page with a custom message
		http.Redirect(w, r, "/login?message=Please+log+in+to+create+a+post.", http.StatusSeeOther)
		return
	}

	// Get the user ID from the session token
	userID, err := repository.GetUserIDBySession(cookie.Value)
	if err != nil || userID == 0 {
		// Redirect to login page with a custom message
		http.Redirect(w, r, "/login?message=Please+log+in+to+create+a+post.", http.StatusSeeOther)
		return
	}

	// Handle GET request to show the post creation form
	if r.Method == http.MethodGet {
		//Get the user object
		user, err := repository.GetUserByID(userID)
		if err != nil {
			log.Println("Error fetching user by ID:", err)
			utils.RenderServerErrorPage(w)
			return
		}
		// Render the create_post.html form
		tmpl, err := template.ParseFiles("web/templates/create_post.html", "web/templates/partials/navbar.html")
		if err != nil {
			log.Println("Could not load template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}
		//Pass the user's profile picture and login status to the template
		data := map[string]interface{}{
			"isLoggedIn":     true,
			"ProfilePicture": user.ProfilePicture,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println("Error executing template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
		}
		return
	}

	// Handle POST request to create a new post
	if r.Method == http.MethodPost {
		// Parse multipart form data to handle file uploads
		err := r.ParseMultipartForm(10 << 20) // Limit to ~10MB
		if err != nil {
			log.Println("Error parsing multipart form:", err)
			utils.RenderServerErrorPage(w)
			return
		}

		// Get form values
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")

		// Get user again
		user, err := repository.GetUserByID(userID)
		if err != nil {
			log.Println("Could not fetch user:", err)
			utils.RenderServerErrorPage(w)
			return
		}

		// Handle image upload
		var imageFilename string
		file, header, err := r.FormFile("image")
		if err == nil {
			defer file.Close()

			imageFilename = header.Filename
			dest, err := utils.SaveUploadedFile(file, imageFilename, "web/static/uploads")
			if err != nil {
				log.Println("Error saving uploaded image:", err)
				imageFilename = "" // fallback to no image
			} else {
				imageFilename = dest // saved path or filename
			}
		} else {
			log.Println("No image uploaded:", err)
		}

		// Validate title and content
		if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
			tmpl, _ := template.ParseFiles("web/templates/create_post.html", "web/templates/partials/navbar.html")
			data := map[string]interface{}{
				"Error":          "Post title and content cannot be empty or spaces only.",
				"Title":          title,
				"Content":        content,
				"Category":       category,
				"ProfilePicture": user.ProfilePicture,
				"isLoggedIn":     true,
			}
			tmpl.Execute(w, data)
			return
		}

		// Save post with image
		err = repository.CreatePost(userID, title, content, category, imageFilename)
		if err != nil {
			log.Println("Error creating post:", err)
			utils.RenderServerErrorPage(w)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// If the method is not allowed, return a Method Not Allowed error
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
