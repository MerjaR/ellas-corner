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

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.FormValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("EditPostHandler: Invalid post ID:", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	sessionUser, err := utils.GetSessionUser(r)
	if err != nil {
		log.Println("EditPostHandler: User not logged in")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	post, err := repository.GetPostByID(postIDStr, sessionUser.ID)
	if err != nil {
		log.Println("EditPostHandler: Error fetching post:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	categories, err := repository.FetchCategories()
	if err != nil {
		log.Println("EditPostHandler: Error fetching categories:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	if r.Method == http.MethodGet {
		isLoggedIn := true
		profilePicture := sessionUser.ProfilePicture

		tmpl, err := template.ParseFiles("web/templates/edit_post.html", "web/templates/partials/navbar.html")
		if err != nil {
			log.Println("EditPostHandler: Error parsing template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		data := viewmodels.EditPostPageData{
			IsLoggedIn:     isLoggedIn,
			ProfilePicture: profilePicture,
			Post:           *post,
			Categories:     categories,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println("EditPostHandler: Error executing template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
		}
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Println("EditPostHandler: Error parsing form:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")
		isDonation := r.FormValue("is_donation") != "on"

		var imagePath string
		file, header, err := r.FormFile("image")
		if err == nil && header.Size > 0 {
			defer file.Close()
			imagePath, err = repository.SaveImageFile(file, header)
			if err != nil {
				log.Println("EditPostHandler: Failed to save image:", err)
				w.WriteHeader(http.StatusInternalServerError)
				utils.RenderServerErrorPage(w)
				return
			}
		} else {
			imagePath = post.Image
		}

		err = repository.UpdatePostWithImage(postID, title, content, category, isDonation, imagePath)
		if err != nil {
			log.Println("EditPostHandler: Error updating post:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	// Unsupported method
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
