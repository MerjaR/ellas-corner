package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils" // Import the utils package for the custom error page
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ProfileHandler: Request received")

	// Check if the user is logged in by checking the session
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("ProfileHandler: No session token found, redirecting to login")
		// Redirect to the login page with a custom message
		http.Redirect(w, r, "/login?message=Please+log+in+to+view+your+profile.", http.StatusSeeOther)
		return
	}

	// Fetch user ID from session
	userID, err := repository.GetUserIDBySession(sessionCookie.Value)
	if err != nil || userID == 0 {
		log.Println("ProfileHandler: Invalid session or user ID not found, redirecting to login")
		// Redirect to the login page with a custom message
		http.Redirect(w, r, "/login?message=Please+log+in+to+view+your+profile.", http.StatusSeeOther)
		return
	}

	log.Printf("ProfileHandler: Logged in user ID: %d\n", userID)

	// Fetch user details (username, email, profile picture)
	user, err := repository.GetUserByID(userID)
	if err != nil {
		log.Println("ProfileHandler: Error fetching user data:", err)
		w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	// Fetch the user's posts, comments, liked posts, and disliked posts
	posts, err := repository.FetchPostsByUser(userID)
	if err != nil {
		log.Println("ProfileHandler: Error fetching user's posts:", err)
		w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	comments, err := repository.FetchCommentsByUser(userID)
	if err != nil {
		log.Println("ProfileHandler: Error fetching user's comments:", err)
		w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	likedPosts, err := repository.FetchLikedPostsByUser(userID)
	if err != nil {
		log.Println("ProfileHandler: Error fetching liked posts:", err)
		w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	dislikedPosts, err := repository.FetchDislikedPostsByUser(userID)
	if err != nil {
		log.Println("ProfileHandler: Error fetching disliked posts:", err)
		w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	log.Println("ProfileHandler: Successfully fetched all data")
	funcMap := template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict expects even number of arguments")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}

	tmpl, err := template.New("profile.html").Funcs(funcMap).ParseFiles(
		"web/templates/profile.html",
		"web/templates/partials/navbar.html",
		"web/templates/partials/post.html",
	)
	if err != nil {
		log.Println("ProfileHandler: Error loading template:", err)
		w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	log.Println("ProfileHandler: Template loaded successfully")

	// Data to pass to the template
	data := map[string]interface{}{
		"Username":                   user.Username,
		"Email":                      user.Email,
		"ProfilePicture":             user.ProfilePicture,
		"Posts":                      posts,         // Posts the user made
		"Comments":                   comments,      // Comments the user made
		"LikedPosts":                 likedPosts,    // Posts the user liked
		"DislikedPosts":              dislikedPosts, // Posts the user disliked
		"isLoggedIn":                 true,          // Indicate that the user is logged in
		"Country":                    user.Country,
		"ShowDonationsInCountryOnly": user.ShowDonationsInCountryOnly,
	}

	// Execute the template
	err = tmpl.ExecuteTemplate(w, "profile.html", data)
	if err != nil {
		log.Println("ProfileHandler: Error executing template:", err)
		w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	log.Println("ProfileHandler: Template executed successfully")
}

func UpdateProfileSettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := repository.GetUserIDBySession(sessionCookie.Value)
	if err != nil || userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse the submitted form values
	country := r.FormValue("country")
	showDonations := r.FormValue("show_donations_in_country_only") == "on"

	// Update user preferences in the DB
	err = repository.UpdateUserPreferences(userID, country, showDonations)
	if err != nil {
		log.Println("UpdateProfileSettingsHandler: Failed to update preferences:", err)
		http.Error(w, "Could not save preferences", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
