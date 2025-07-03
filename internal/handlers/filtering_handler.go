package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"ellas-corner/internal/viewmodels"
	"html/template"
	"log"
	"net/http"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("FilterHandler: Request received")

	// Check session to determine login state
	sessionUser, err := utils.GetSessionUser(r)
	isLoggedIn := err == nil
	userID := 0
	profilePicture := ""
	if isLoggedIn {
		userID = sessionUser.ID
		profilePicture = sessionUser.ProfilePicture
	}

	// Read query parameters used for filtering
	category := r.URL.Query().Get("category")
	createdPosts := r.URL.Query().Get("created_posts")
	likedPosts := r.URL.Query().Get("liked_posts")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Fetch posts based on filters
	posts, err := repository.FetchFilteredPosts(category, createdPosts, likedPosts, startDate, endDate, userID, isLoggedIn)
	if err != nil {
		log.Println("FilterHandler: Error fetching filtered posts:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	// Fetch available categories for filter options
	categories, err := repository.FetchCategories()
	if err != nil {
		log.Println("FilterHandler: Error fetching categories:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	// Load and parse template files
	tmpl, err := template.ParseFiles(
		"web/templates/filter_results.html",
		"web/templates/partials/navbar.html",
		"web/templates/partials/post.html",
	)
	if err != nil {
		log.Println("FilterHandler: Error parsing template:", err)
		utils.RenderServerErrorPage(w)
		return
	}

	// Prepare data for template rendering
	data := viewmodels.FilterPageData{
		IsLoggedIn:     isLoggedIn,
		ProfilePicture: profilePicture,
		Posts:          posts,
		Categories:     categories,
		Category:       category,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("FilterHandler: Error executing template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	log.Println("FilterHandler: Template executed successfully")
}
