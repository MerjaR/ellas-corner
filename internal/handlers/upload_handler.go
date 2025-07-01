package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils" // Import utils package for custom error page
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// UploadProfilePictureHandler handles the profile picture upload
func UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	sessionUser, err := utils.GetSessionUser(r)
	if err != nil {
		http.Error(w, "Please log in to upload a profile picture", http.StatusUnauthorized)
		return
	}

	userID := sessionUser.ID

	// Parse the form to retrieve the file
	err = r.ParseMultipartForm(10 << 20) // Max 10MB file size
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("profile_picture")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate the file type (only allow image files)
	ext := filepath.Ext(handler.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		http.Error(w, "Only image files are allowed (jpg, png, gif)", http.StatusBadRequest)
		return
	}

	// Create a file on the server to store the uploaded image
	filePath := fmt.Sprintf("web/static/profile_pictures/user_%d%s", userID, ext)
	out, err := os.Create(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Set 500 status code
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}
	defer out.Close()

	// Copy the file to the server
	_, err = io.Copy(out, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Set 500 status code
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	// Update the user's profile picture in the database
	err = repository.UpdateProfilePicture(userID, filepath.Base(filePath))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Set 500 status code
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	// Redirect back to the profile page
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
