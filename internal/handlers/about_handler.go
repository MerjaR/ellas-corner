package handlers

import (
	"ellas-corner/internal/repository"
	"html/template"
	"log"
	"net/http"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/about.html", "web/templates/partials/navbar.html")
	if err != nil {
		http.Error(w, "Error loading About page", http.StatusInternalServerError)
		return
	}

	var isLoggedIn bool
	var profilePicture string

	// Check for session cookie
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

	data := map[string]interface{}{
		"isLoggedIn":     isLoggedIn,
		"ProfilePicture": profilePicture,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Error executing About page template:", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
