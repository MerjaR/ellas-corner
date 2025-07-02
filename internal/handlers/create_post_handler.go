package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"ellas-corner/internal/viewmodels"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in (session token)
	sessionUser, err := utils.GetSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login?message=Please+log+in+to+create+a+post.", http.StatusSeeOther)
		return
	}

	// Handle GET request to show the post creation form
	if r.Method == http.MethodGet {

		// Render the create_post.html form
		tmpl, err := template.ParseFiles("web/templates/create_post.html", "web/templates/partials/navbar.html")
		if err != nil {
			log.Println("Could not load template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}
		//Pass the user's profile picture and login status to the template
		data := viewmodels.CreatePostPageData{
			IsLoggedIn:     true,
			ProfilePicture: sessionUser.ProfilePicture,
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
		user, err := repository.GetUserByID(sessionUser.ID)
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
				imageFilename = "placeholder.jpg"
			} else {
				imageFilename = dest
			}
		} else {
			log.Println("No image uploaded; using placeholder.")
			imageFilename = "placeholder.jpg"
		}

		isDonation := r.FormValue("is_donation") == "on"

		// Validate title and content
		if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
			tmpl, _ := template.ParseFiles("web/templates/create_post.html", "web/templates/partials/navbar.html")
			data := viewmodels.CreatePostPageData{
				Error:          "Post title and content cannot be empty or spaces only.",
				Title:          title,
				Content:        content,
				Category:       category,
				ProfilePicture: user.ProfilePicture,
				IsLoggedIn:     true,
			}

			tmpl.Execute(w, data)
			return
		}

		// Save post with image
		err = repository.CreatePost(sessionUser.ID, title, content, category, imageFilename, isDonation, user.Country)

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
