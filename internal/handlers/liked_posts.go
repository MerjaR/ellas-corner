package handlers

import (
	"ellas-corner/internal/db"
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

	user, err := repository.GetUserByID(userID)
	if err != nil {
		utils.RenderServerErrorPage(w)
		return
	}

	likedPosts, err := repository.FetchLikedPostsByUser(userID)
	if err != nil {
		utils.RenderServerErrorPage(w)
		return
	}

	selectedCategories := r.URL.Query()["category"]

	var curatedItems []db.BabyBoxItem
	for _, cat := range selectedCategories {
		if items, ok := db.CuratedBabyBox[cat]; ok {
			curatedItems = append(curatedItems, items...)
		}
	}

	tmpl, err := template.ParseFiles("web/templates/liked_posts.html", "web/templates/partials/navbar.html")
	if err != nil {
		utils.RenderServerErrorPage(w)
		return
	}

	data := map[string]interface{}{
		"LikedPosts":     likedPosts,
		"isLoggedIn":     true,
		"ProfilePicture": user.ProfilePicture,
		"Username":       user.Username,
		"CuratedItems":   curatedItems,
	}

	tmpl.Execute(w, data)
}
