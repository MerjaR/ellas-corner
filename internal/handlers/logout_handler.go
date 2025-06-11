package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"log"
	"net/http"
	"time"
)

// LogoutHandler logs the user out by clearing the session cookie and removing the session from the database
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Delete session from the database
		err = repository.DeleteSession(cookie.Value)
		if err != nil {
			log.Println("Error deleting session from database:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w) // Show custom 5xx error page if deleting session fails
			return
		}
	}

	// Clear the session cookie by setting its expiration to the past
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",             // Make sure to clear the cookie across the entire site
		Expires: time.Unix(0, 0), // Set expiration to a date in the past
	})

	// Redirect back to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
