package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"ellas-corner/internal/viewmodels"
	"html/template"
	"log"
	"net/http"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SearchHandler: Request received")

	// Get the search query
	searchQuery := r.URL.Query().Get("q")
	if searchQuery == "" {
		http.Error(w, "Please provide a search query", http.StatusBadRequest)
		return
	}

	sessionUser, err := utils.GetSessionUser(r)
	isLoggedIn := false
	var userID int
	var profilePicture string
	var currentUser repository.User // needed for donation logic

	if err == nil {
		isLoggedIn = true
		userID = sessionUser.ID
		profilePicture = sessionUser.ProfilePicture

		currentUser, err = repository.GetUserByID(userID)
		if err != nil {
			log.Println("SearchHandler: Error fetching full user profile:", err)
			utils.RenderServerErrorPage(w)
			return
		}
	}

	// Fetch posts matching the query
	posts, err := repository.SearchPosts(searchQuery, userID, isLoggedIn)
	if err != nil {
		log.Println("Error searching posts:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	for i := range posts {
		if posts[i].IsDonation {
			if isLoggedIn && currentUser.ShowDonationsInCountryOnly {
				posts[i].ShowDonatedLabel = posts[i].DonationCountry == currentUser.Country
			} else {
				posts[i].ShowDonatedLabel = true
			}
		}
	}

	for i := range posts {
		comments, err := repository.FetchCommentsForPost(posts[i].ID, userID)
		if err != nil {
			log.Println("Error fetching comments for post ID", posts[i].ID, ":", err)
			continue // skip adding comments if there's an error
		}
		posts[i].Comments = comments
	}

	// Parse the search results template and navbar
	tmpl, err := template.ParseFiles("web/templates/search_results.html", "web/templates/partials/navbar.html", "web/templates/partials/post.html")
	if err != nil {
		log.Println("SearchHandler: Error parsing template", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	// Prepare data for the template
	data := viewmodels.SearchPageData{
		IsLoggedIn:             isLoggedIn,
		ProfilePicture:         profilePicture,
		SearchQuery:            searchQuery,
		Posts:                  posts,
		ShowEditControls:       false,
		ShowCommentFormForPost: 0,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("SearchHandler: Error executing template", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	log.Println("SearchHandler: Template executed successfully")
}
