package utils

import (
	"html/template"
	"log"
	"net/http"
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
