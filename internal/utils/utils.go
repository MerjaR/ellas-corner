package utils

import (
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// renderServerErrorPage renders a custom 500 error page
func RenderServerErrorPage(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError) // Set HTTP status to 500

	tmpl, err := template.ParseFiles("web/templates/500.html")
	if err != nil {
		log.Println("Error loading 500 error template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("Error executing 500 error template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func SaveUploadedFile(file multipart.File, filename, uploadPath string) (string, error) {
	// Make sure the directory exists
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	dstPath := filepath.Join(uploadPath, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return filename, nil
}


