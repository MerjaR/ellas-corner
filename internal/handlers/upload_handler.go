package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// UploadProfilePictureHandler handles the profile picture upload
func UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	const maxUploadSize = 10 << 20 // 10MB

	sessionUser, err := utils.GetSessionUser(r)
	if err != nil {
		http.Error(w, "Please log in to upload a profile picture", http.StatusUnauthorized)
		return
	}
	userID := sessionUser.ID

	// Enforce file size limit before reading body
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// Parse the multipart form
	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		log.Println("UploadProfilePictureHandler: Error parsing multipart form:", err)
		http.Error(w, "File too large or invalid", http.StatusBadRequest)
		return
	}

	// Retrieve file
	file, _, err := r.FormFile("profile_picture")
	if err != nil {
		log.Println("UploadProfilePictureHandler: Error retrieving file:", err)
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check MIME type by reading first 512 bytes
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		log.Println("UploadProfilePictureHandler: Error reading file:", err)
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	filetype := http.DetectContentType(buffer)
	_, err = file.Seek(0, 0) // reset file pointer
	if err != nil {
		log.Println("UploadProfilePictureHandler: Error resetting file pointer:", err)
		http.Error(w, "Failed to process file", http.StatusInternalServerError)
		return
	}

	// Allow only valid image types
	allowedTypes := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
	}
	ext, ok := allowedTypes[filetype]
	if !ok {
		log.Printf("UploadProfilePictureHandler: Rejected file type: %s", filetype)
		http.Error(w, "Only image files are allowed (jpg, png, gif)", http.StatusBadRequest)
		return
	}

	// Generate unique filename with timestamp to avoid overwrite
	filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join("web/static/profile_pictures", filename)

	// Create destination file
	out, err := os.Create(filePath)
	if err != nil {
		log.Println("UploadProfilePictureHandler: Error creating file:", err)
		utils.RenderServerErrorPage(w)
		return
	}
	defer out.Close()

	// Copy uploaded data to server file
	_, err = io.Copy(out, file)
	if err != nil {
		log.Println("UploadProfilePictureHandler: Error saving file:", err)
		utils.RenderServerErrorPage(w)
		return
	}

	// Update user record with new filename
	err = repository.UpdateProfilePicture(userID, filename)
	if err != nil {
		log.Println("UploadProfilePictureHandler: Failed to update user profile:", err)
		utils.RenderServerErrorPage(w)
		return
	}

	// Redirect to profile
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
