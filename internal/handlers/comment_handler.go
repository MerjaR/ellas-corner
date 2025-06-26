package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionUser, err := utils.GetSessionUser(r)
	if err != nil {
		http.Error(w, "Unauthorized. Please log in to comment.", http.StatusUnauthorized)
		return
	}

	// Get form values (post_id, content)
	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	// Check if the comment is empty or contains only whitespace
	if strings.TrimSpace(content) == "" {
		// Fetch the post again, as well as its comments, and show the error in the template
		post, err := repository.GetPostByID(postID)
		if err != nil {
			log.Println("Error fetching post:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}
		comments, err := repository.FetchCommentsForPost(post.ID)
		if err != nil {
			log.Println("Error fetching comments:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}
		post.Comments = comments

		// Render the template with the error message
		tmpl, err := template.ParseFiles("web/templates/index.html")
		if err != nil {
			log.Println("Could not load template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		// Pass the error message and comment content to the template
		data := map[string]interface{}{
			"Posts":                  []repository.Post{*post},
			"ErrorMessage":           "Comment cannot be empty or only spaces",
			"CommentContent":         content, // Preserve user input
			"ShowCommentFormForPost": post.ID, // Show form again for this post
			"IsLoggedIn":             true,    // Make sure to pass this correctly
		}

		// Show the error on the same post page
		if err := tmpl.Execute(w, data); err != nil {
			log.Println("Error executing template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
		}
		return
	}

	// Save the comment to the database
	err = repository.CreateComment(sessionUser.ID, postID, content)
	if err != nil {
		log.Println("Error creating comment:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	// Redirect back to the post after the comment is added
	http.Redirect(w, r, "/?showCommentFormForPost="+postID, http.StatusSeeOther)
}
