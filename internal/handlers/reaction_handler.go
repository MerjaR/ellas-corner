package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"ellas-corner/internal/viewmodels"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ReactionHandler: Request received")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := 0
	sessionUser, err := utils.GetSessionUser(r)
	if err != nil {
		log.Println("User not logged in or invalid session")
		renderHomeWithError(w, "You must be logged in to react.", userID)
		return
	}
	userID = sessionUser.ID

	// Get form values
	postIDStr := r.FormValue("post_id")
	reaction := r.FormValue("reaction")
	log.Printf("ReactionHandler: Post ID: %s, Reaction: %s", postIDStr, reaction)

	// Convert postIDStr to int
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("Invalid post ID:", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Add the reaction to the database
	err = repository.AddReaction(userID, postID, reaction)
	if err != nil {
		log.Println("Error adding reaction:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	// Redirect back to the home page after reacting
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Helper function to render the home page with an error message (without r *http.Request)
func renderHomeWithError(w http.ResponseWriter, errorMessage string, userID int) {
	posts, err := repository.FetchPosts(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	tmpl, err := template.ParseFiles(
		"web/templates/index.html",
		"web/templates/partials/post.html",
		"web/templates/partials/navbar.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	data := viewmodels.HomePageData{
		IsLoggedIn:   false,
		Posts:        posts,
		ErrorMessage: errorMessage,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error rendering index with error message:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
	}
}

func CommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("CommentReactionHandler: Request received")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionUser, err := utils.GetSessionUser(r)
	userID := 0
	if err == nil {
		userID = sessionUser.ID
	} else {
		log.Println("User not logged in or invalid session")
		renderHomeWithError(w, "You must be logged in to react.", userID)
		return
	}

	// Get form values
	commentIDStr := r.FormValue("comment_id")
	reaction := r.FormValue("reaction")
	log.Printf("CommentReactionHandler: Comment ID: %s, Reaction: %s", commentIDStr, reaction)

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		log.Println("Invalid comment ID:", err)
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Add the reaction to the comment in the database
	err = repository.AddCommentReaction(userID, commentID, reaction)
	if err != nil {
		log.Println("Error adding reaction to comment:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	// Redirect back to the home page after reacting
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
