package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils" // Import the utils package for the custom error page
	"html/template"
	"log"
	"net/http"
)

// SearchHandler handles the search functionality
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SearchHandler: Request received")

	// Get the search query from the form
	searchQuery := r.URL.Query().Get("q")
	if searchQuery == "" {
		http.Error(w, "Please provide a search query", http.StatusBadRequest)
		return
	}

	// Fetch the posts that match the search query
	posts, err := repository.SearchPosts(searchQuery)
	if err != nil {
		log.Println("Error searching posts:", err)
		w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	// Parse the search results template
	tmpl, err := template.ParseFiles("web/templates/search_results.html")
	if err != nil {
		log.Println("SearchHandler: Error parsing template", err)
		w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	// Render the template with the search results
	data := map[string]interface{}{
		"SearchQuery": searchQuery,
		"Posts":       posts, // The search results (posts)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("SearchHandler: Error executing template", err)
		w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	log.Println("SearchHandler: Template executed successfully")
}
