package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"html/template"
	"log"
	"net/http"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("FilteringHandler: Request received")

	// Check if the user is logged in
	sessionCookie, err := r.Cookie("session_token")
	isLoggedIn := false
	var userID int

	if err == nil {
		// Validate session token and check if user is logged in
		userID, err = repository.GetUserIDBySession(sessionCookie.Value)
		if err == nil && userID != 0 {
			isLoggedIn = true
		}
	}

	// If user is logged in, fetch profile picture
	var ProfilePicture string
	if isLoggedIn {
		user, err := repository.GetUserByID(userID)
		if err == nil {
			ProfilePicture = user.ProfilePicture
		}
	}

	// Capture filter parameters from the query string
	category := r.URL.Query().Get("category")
	createdPosts := r.URL.Query().Get("created_posts")
	likedPosts := r.URL.Query().Get("liked_posts")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Fetch filtered posts based on the filter parameters
	posts, err := repository.FetchFilteredPosts(category, createdPosts, likedPosts, startDate, endDate, userID, isLoggedIn)
	if err != nil {
		log.Println("Error fetching filtered posts:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	// Fetch the available categories
	categories, err := repository.FetchCategories()
	if err != nil {
		log.Println("Error fetching categories:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}
	log.Println("Categories fetched:", categories) // Log categories fetched

	// Parse the filtered posts template
	tmpl, err := template.ParseFiles("web/templates/filter_results.html", "web/templates/partials/navbar.html", "web/templates/partials/post.html")
	if err != nil {
		log.Println("FilteringHandler: Error parsing template", err)
		utils.RenderServerErrorPage(w)
		return
	}

	// Render the template with filtered posts and available categories
	data := map[string]interface{}{
		"isLoggedIn":     isLoggedIn,
		"Posts":          posts,
		"Categories":     categories, // Pass categories to the template
		"ProfilePicture": ProfilePicture,
		"Category":       category,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("FilteringHandler: Error executing template", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	log.Println("FilteringHandler: Template executed successfully")
}
