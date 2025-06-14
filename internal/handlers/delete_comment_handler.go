package handlers

import (
	"ellas-corner/internal/repository"
	"log"
	"net/http"
	"strconv"
)

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	commentIDStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		log.Println("Invalid comment ID:", err)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	err = repository.DeleteComment(commentID)
	if err != nil {
		log.Println("Error deleting comment:", err)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
