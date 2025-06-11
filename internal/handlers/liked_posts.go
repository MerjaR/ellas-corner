package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"html/template"
	"log"
	"net/http"
)

func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Get session
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := repository.GetUserIDBySession(sessionCookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch liked posts
	posts, err := repository.FetchLikedPosts(userID)
	if err != nil {
		log.Println("Error fetching liked posts:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/liked_posts.html")
	if err != nil {
		log.Println("Error parsing liked posts template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	data := map[string]interface{}{
		"IsLoggedIn": true,
		"Posts":      posts,
		"PageTitle":  "Liked Posts",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error rendering liked posts:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
	}
}
