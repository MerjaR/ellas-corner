package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// RegisterHandler handles the registration form and logic
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/register.html", "web/templates/partials/navbar_register.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w) // Error handling
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// Handle form submission
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Encrypt the password
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w) // Error handling
			return
		}

		// Randomly assign a profile picture (no need to seed)
		pictureOptions := []string{"1.png", "2.png", "3.png"}
		randomPicture := pictureOptions[rand.Intn(len(pictureOptions))]

		// Save the user to the database with the random profile picture
		err = repository.CreateUser(username, email, hashedPassword, randomPicture)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				// Render the registration form again with an error message
				tmpl, _ := template.ParseFiles("web/templates/register.html", "web/templates/partials/navbar_register.html")
				tmpl.Execute(w, map[string]interface{}{
					"Error": "Email or username already in use. Please try a different one.",
				})
			} else {
				log.Println("Error creating user:", err)
				w.WriteHeader(http.StatusInternalServerError)
				utils.RenderServerErrorPage(w) // Error handling
			}
			return
		}

		// Redirect to login page with a success message
		http.Redirect(w, r, "/login?message=Thank+you+for+joining+Ella's+Corner!+Your+registration+was+successful,+please+log+in.", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// LoginHandler handles user login logic
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Check for a message in the query parameters
		message := r.URL.Query().Get("message")

		// Prepare the data for the template
		data := map[string]interface{}{}
		if message != "" {
			data["Message"] = message
		}

		// Render the login template with any message
		tmpl, err := template.ParseFiles("web/templates/login.html", "web/templates/partials/navbar_minimal.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w) // Error handling
			return
		}
		tmpl.Execute(w, data)
	} else if r.Method == http.MethodPost {
		// Handle form submission (login logic)
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Fetch user from database by email
		user, err := repository.GetUserByEmail(email)
		if err != nil || user == nil {
			tmpl, _ := template.ParseFiles("web/templates/login.html", "web/templates/partials/navbar_minimal.html")
			tmpl.Execute(w, map[string]interface{}{"Error": "Invalid email or password"})
			return
		}

		// Check if the password matches the hashed password
		if !utils.CheckPasswordHash(password, user.Password) {
			tmpl, _ := template.ParseFiles("web/templates/login.html", "web/templates/partials/navbar_minimal.html")
			tmpl.Execute(w, map[string]interface{}{"Error": "Invalid email or password"})
			return
		}

		// Generate a session token
		sessionToken := utils.GenerateSessionToken()

		// Save the session token to the sessions table
		err = repository.SaveSessionToken(user.ID, sessionToken)
		if err != nil {
			log.Println("Error saving session token:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w) // Error handling
			return
		}

		// Set the session token in the cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(24 * time.Hour), // Cookie expires in 24 hours
		})

		// Redirect to the homepage after successful login
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// AcceptCookiesHandler saves the user's cookie consent in the database or browser cookies
func AcceptCookiesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AcceptCookiesHandler: Request received")

	// Check if the user is logged in by checking for a valid session associated with a user
	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		userID, err := repository.GetUserIDBySession(sessionCookie.Value)
		if err == nil && userID != 0 {
			// If a valid session is found for a user, save the cookie consent in the database
			err = repository.SaveCookieConsent(userID, true)
			if err != nil {
				log.Println("Error saving cookie consent:", err)
				w.WriteHeader(http.StatusInternalServerError)
				utils.RenderServerErrorPage(w) // Error handling
				return
			}
		} else {
			// If the session token is not associated with a user, treat as non-logged-in
			http.SetCookie(w, &http.Cookie{
				Name:    "consent_given",
				Value:   "true",
				Expires: time.Now().Add(365 * 24 * time.Hour), // 1-year expiration
				Path:    "/",                                  // Site-wide cookie
			})
			log.Println("Cookie set for non-logged-in user")
		}
	} else {
		// If no session token is found at all, treat as non-logged-in
		http.SetCookie(w, &http.Cookie{
			Name:    "consent_given",
			Value:   "true",
			Expires: time.Now().Add(365 * 24 * time.Hour), // 1-year expiration
			Path:    "/",                                  // Site-wide cookie
		})
		log.Println("Cookie set for non-logged-in user (no session token)")
	}

	// Redirect back to the homepage after consent is saved
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
