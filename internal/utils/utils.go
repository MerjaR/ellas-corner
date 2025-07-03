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

// RenderServerErrorPage renders a custom 500 error page.
// It ensures WriteHeader is only called once.
func RenderServerErrorPage(w http.ResponseWriter) {
	// Check if header was already written
	if w.Header().Get("Content-Type") == "" {
		w.WriteHeader(http.StatusInternalServerError)
	}

	tmpl, err := template.ParseFiles("web/templates/500.html")
	if err != nil {
		log.Println("RenderServerErrorPage: error loading 500.html:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println("RenderServerErrorPage: error executing 500 template:", err)
		// Only call WriteHeader if response hasn’t started
		if w.Header().Get("Content-Type") == "" {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
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
