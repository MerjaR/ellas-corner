package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"ellas-corner/internal/viewmodels"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// AddCommentHandler processes user-submitted comments.
// Requires user to be logged in and content to be non-empty.
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

	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	if strings.TrimSpace(content) == "" {
		post, err := repository.GetPostByID(postID, sessionUser.ID)
		if err != nil || post == nil {
			log.Println("AddCommentHandler: Error fetching post:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		comments, err := repository.FetchCommentsForPost(post.ID, sessionUser.ID)
		if err == nil {
			post.Comments = comments
		}

		const indexTemplate = "web/templates/index.html"
		tmpl, err := template.ParseFiles(indexTemplate)
		if err != nil {
			log.Println("AddCommentHandler: Error loading template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		data := viewmodels.HomePageData{
			Posts:                  []repository.Post{*post},
			ErrorMessage:           "Comment cannot be empty or only spaces",
			ShowCommentFormForPost: post.ID,
			IsLoggedIn:             true,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println("AddCommentHandler: Error executing template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
		}
		return
	}

	err = repository.CreateComment(sessionUser.ID, postID, content)
	if err != nil {
		log.Println("AddCommentHandler: Error creating comment:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	http.Redirect(w, r, "/?showCommentFormForPost="+postID, http.StatusSeeOther)
}
