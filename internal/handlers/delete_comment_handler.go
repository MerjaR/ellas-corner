package handlers

import (
	"ellas-corner/internal/repository"
	"log"
	"net/http"
	"strconv"
)

// DeleteCommentHandler handles the deletion of a user comment from their profile page.
func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse comment ID from form data
	commentIDStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		log.Println("DeleteCommentHandler: Invalid comment ID:", err)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	// Attempt to delete the comment
	if err := repository.DeleteComment(commentID); err != nil {
		log.Println("DeleteCommentHandler: Error deleting comment:", err)
	} else {
		log.Printf("DeleteCommentHandler: Comment %d deleted successfully\n", commentID)
	}

	// Redirect user back to their profile page
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
