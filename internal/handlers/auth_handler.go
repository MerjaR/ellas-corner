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

// RegisterHandler renders the registration form on GET,
// and handles user registration on POST.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	const registerTemplate = "web/templates/register.html"
	const navbarTemplate = "web/templates/partials/navbar_register.html"

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(registerTemplate, navbarTemplate)
		if err != nil {
			log.Println("RegisterHandler: Error parsing template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if username == "" || email == "" || password == "" {
			tmpl, err := template.ParseFiles(registerTemplate, navbarTemplate)
			if err == nil {
				tmpl.Execute(w, map[string]interface{}{
					"Error": "Please fill in all fields.",
				})
			}
			return
		}

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		pictureOptions := []string{"1.png", "2.png", "3.png"}
		randomPicture := pictureOptions[rand.Intn(len(pictureOptions))]

		err = repository.CreateUser(username, email, hashedPassword, randomPicture)
		if err != nil {
			log.Println("Error creating user:", err)
			tmpl, tmplErr := template.ParseFiles(registerTemplate, navbarTemplate)
			if tmplErr == nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
				tmpl.Execute(w, map[string]interface{}{
					"Error": "Email or username already in use. Please try a different one.",
				})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				utils.RenderServerErrorPage(w)
			}
			return
		}

		http.Redirect(w, r, "/login?message=Thank+you+for+joining+Ella's+Corner!+Your+registration+was+successful,+please+log+in.", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// LoginHandler renders the login form on GET,
// and handles authentication on POST.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	const loginTemplate = "web/templates/login.html"
	const navbarTemplate = "web/templates/partials/navbar_minimal.html"

	switch r.Method {
	case http.MethodGet:
		message := r.URL.Query().Get("message")
		data := map[string]interface{}{}
		if message != "" {
			data["Message"] = message
		}

		tmpl, err := template.ParseFiles(loginTemplate, navbarTemplate)
		if err != nil {
			log.Println("LoginHandler: Error parsing template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println("LoginHandler: Error executing template:", err)
		}

	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := repository.GetUserByEmail(email)
		if err != nil || user == nil {
			log.Println("LoginHandler: Invalid email or user not found")
			renderLoginError(w, "Invalid email or password")
			return
		}

		if !utils.CheckPasswordHash(password, user.Password) {
			log.Println("LoginHandler: Incorrect password for user:", user.Email)
			renderLoginError(w, "Invalid email or password")
			return
		}

		sessionToken := utils.GenerateSessionToken()
		if err := repository.SaveSessionToken(user.ID, sessionToken); err != nil {
			log.Println("LoginHandler: Error saving session token:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Helper to render login template with an error
func renderLoginError(w http.ResponseWriter, errorMsg string) {
	tmpl, err := template.ParseFiles("web/templates/login.html", "web/templates/partials/navbar_minimal.html")
	if err != nil {
		log.Println("renderLoginError: Error loading login template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	if err := tmpl.Execute(w, map[string]interface{}{"Error": errorMsg}); err != nil {
		log.Println("renderLoginError: Error executing login template:", err)
	}
}

// AcceptCookiesHandler records user cookie consent,
// either in the DB (if logged in) or as a browser cookie.
func AcceptCookiesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AcceptCookiesHandler: Request received")

	cookie := &http.Cookie{
		Name:    "consent_given",
		Value:   "true",
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Path:    "/",
	}

	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		userID, err := repository.GetUserIDBySession(sessionCookie.Value)
		if err == nil && userID != 0 {
			if err := repository.SaveCookieConsent(userID, true); err != nil {
				log.Println("AcceptCookiesHandler: Error saving consent to DB:", err)
				w.WriteHeader(http.StatusInternalServerError)
				utils.RenderServerErrorPage(w)
				return
			}
			log.Println("Consent saved for logged-in user:", userID)
		} else {
			http.SetCookie(w, cookie)
			log.Println("Consent cookie set for unauthenticated session token")
		}
	} else {
		http.SetCookie(w, cookie)
		log.Println("Consent cookie set for user with no session token")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
