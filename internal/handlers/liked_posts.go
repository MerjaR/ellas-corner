package handlers

import (
	"ellas-corner/internal/db"
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"ellas-corner/internal/viewmodels"
	"html/template"
	"net/http"
)

func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	sessionUser, err := utils.GetSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	likedPosts, err := repository.FetchLikedPostsByUser(sessionUser.ID)
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

	data := viewmodels.LikedPostsPageData{
		IsLoggedIn:     true,
		ProfilePicture: sessionUser.ProfilePicture,
		Username:       sessionUser.Username,
		LikedPosts:     likedPosts,
		CuratedItems:   curatedItems,
	}

	tmpl.Execute(w, data)
}
