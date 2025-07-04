package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"log"
	"net/http"
	"time"
)

// LogoutHandler logs the user out by clearing the session cookie and deleting the session from the database
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("LogoutHandler: Request received")

	// Only attempt session deletion if cookie is present
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Try deleting session from DB
		if err := repository.DeleteSession(cookie.Value); err != nil {
			log.Println("LogoutHandler: Error deleting session from DB:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}
		log.Println("LogoutHandler: Session deleted from DB")
	} else {
		log.Println("LogoutHandler: No session token found; user may already be logged out")
	}

	// Clear the cookie regardless
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",             // Applies site-wide
		Expires: time.Unix(0, 0), // Expire it immediately
	})

	log.Println("LogoutHandler: Session cookie cleared; redirecting to home")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
