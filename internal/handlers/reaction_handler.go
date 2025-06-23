package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils" // Import the utils package for the custom error page
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// ReactionHandler handles likes and dislikes for posts
func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ReactionHandler: Request received")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is logged in
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("ReactionHandler: User not logged in, rendering index with error")
		renderHomeWithError(w, "You must be logged in to like or dislike a post.")
		return
	}

	// Get user ID from session
	userID, err := repository.GetUserIDBySession(cookie.Value)
	if err != nil || userID == 0 {
		log.Println("Unauthorized: Invalid session or user ID")
		renderHomeWithError(w, "Only registered users can react to posts. Please log in to react.")
		return
	}

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
		w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	// Redirect back to the home page after reacting
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Helper function to render the home page with an error message (without r *http.Request)
func renderHomeWithError(w http.ResponseWriter, errorMessage string) {
	// Fetch the posts
	posts, err := repository.FetchPosts()
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

	// ✅ Don't set w.WriteHeader here — it's just a validation message
	data := map[string]interface{}{
		"Posts":        posts,
		"ErrorMessage": errorMessage,
		"IsLoggedIn":   false,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error rendering index with error message:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
	}
}

// CommentReactionHandler handles likes and dislikes for comments
func CommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("CommentReactionHandler: Request received")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is logged in
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("CommentReactionHandler: User not logged in")
		renderHomeWithError(w, "You must be logged in to like or dislike a comment.")
		return
	}

	// Get user ID from session
	userID, err := repository.GetUserIDBySession(cookie.Value)
	if err != nil || userID == 0 {
		log.Println("Unauthorized: Invalid session or user ID")
		renderHomeWithError(w, "Only registered users can react to comments. Please log in to react.")
		return
	}

	// Get form values
	commentIDStr := r.FormValue("comment_id")
	reaction := r.FormValue("reaction")
	log.Printf("CommentReactionHandler: Comment ID: %s, Reaction: %s", commentIDStr, reaction)

	// Convert commentIDStr to int
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
		w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		utils.RenderServerErrorPage(w)                // Render custom error page
		return
	}

	// Redirect back to the home page after reacting
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
