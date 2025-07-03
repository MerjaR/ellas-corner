package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"log"
	"net/http"
	"strconv"
)

// DeletePostHandler deletes a post based on its ID (requires POST method)
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("DeletePostHandler: Invalid post ID:", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if err := repository.DeletePost(postID); err != nil {
		log.Println("DeletePostHandler: Error deleting post:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	log.Printf("DeletePostHandler: Post %d deleted successfully\n", postID)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
