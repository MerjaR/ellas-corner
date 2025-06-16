package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"html/template"
	"net/http"
)

func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {
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

	likedPosts, err := repository.FetchLikedPostsByUser(userID)
	if err != nil {
		utils.RenderServerErrorPage(w)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/liked_posts.html", "web/templates/partials/navbar.html")
	if err != nil {
		utils.RenderServerErrorPage(w)
		return
	}

	data := map[string]interface{}{
		"LikedPosts": likedPosts,
		"isLoggedIn": true,
	}

	tmpl.Execute(w, data)
}
