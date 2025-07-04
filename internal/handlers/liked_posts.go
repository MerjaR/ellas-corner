package handlers

import (
	"ellas-corner/internal/db"
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"ellas-corner/internal/viewmodels"
	"html/template"
	"log"
	"net/http"
)

// LikedPostsHandler shows posts the user has liked, and curated baby box items
func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	sessionUser, err := utils.GetSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch posts liked by the current user
	likedPosts, err := repository.FetchLikedPostsByUser(sessionUser.ID)
	if err != nil {
		log.Println("LikedPostsHandler: Error fetching liked posts:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	// Handle optional baby box category filtering
	selectedCategories := r.URL.Query()["category"]
	var curatedItems []db.BabyBoxItem
	for _, cat := range selectedCategories {
		if items, ok := db.CuratedBabyBox[cat]; ok {
			curatedItems = append(curatedItems, items...)
		}
	}

	// Parse the template
	tmpl, err := template.ParseFiles(
		"web/templates/liked_posts.html",
		"web/templates/partials/navbar.html",
	)
	if err != nil {
		log.Println("LikedPostsHandler: Error parsing templates:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	// Prepare data for rendering
	data := viewmodels.LikedPostsPageData{
		IsLoggedIn:     true,
		ProfilePicture: sessionUser.ProfilePicture,
		Username:       sessionUser.Username,
		LikedPosts:     likedPosts,
		CuratedItems:   curatedItems,
	}

	// Execute the template
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("LikedPostsHandler: Error executing template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}
}
