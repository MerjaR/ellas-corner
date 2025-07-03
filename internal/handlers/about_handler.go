package handlers

import (
	"ellas-corner/internal/utils"
	"ellas-corner/internal/viewmodels"
	"html/template"
	"log"
	"net/http"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/about.html", "web/templates/partials/navbar.html")
	if err != nil {
		log.Println("Error parsing About page templates:", err)
		http.Error(w, "Error loading About page", http.StatusInternalServerError)
		return
	}

	var isLoggedIn bool
	var profilePicture string

	// Determine login state
	sessionUser, err := utils.GetSessionUser(r)
	if err == nil {
		isLoggedIn = true
		profilePicture = sessionUser.ProfilePicture
	}

	data := viewmodels.HomePageData{
		IsLoggedIn:     isLoggedIn,
		ProfilePicture: profilePicture,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Error rendering About page:", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
