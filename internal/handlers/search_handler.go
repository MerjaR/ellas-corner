package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils" // Import the utils package for the custom error page
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

	if err == nil {
		isLoggedIn = true
		userID = sessionUser.ID
		profilePicture = sessionUser.ProfilePicture
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
		comments, err := repository.FetchCommentsForPost(posts[i].ID)
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
	data := map[string]interface{}{
		"SearchQuery":    searchQuery,
		"Posts":          posts,
		"isLoggedIn":     isLoggedIn,
		"ProfilePicture": profilePicture,
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
