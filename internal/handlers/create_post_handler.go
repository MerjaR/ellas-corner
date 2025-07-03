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

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	const (
		postTemplate   = "web/templates/create_post.html"
		navbarTemplate = "web/templates/partials/navbar.html"
		uploadDir      = "web/static/uploads"
	)

	sessionUser, err := utils.GetSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/login?message=Please+log+in+to+create+a+post.", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles(postTemplate, navbarTemplate)
		if err != nil {
			log.Println("CreatePostHandler: Error parsing template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		data := viewmodels.CreatePostPageData{
			IsLoggedIn:     true,
			ProfilePicture: sessionUser.ProfilePicture,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println("CreatePostHandler: Error executing template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
		}
		return

	case http.MethodPost:
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Println("CreatePostHandler: Error parsing multipart form:", err)
			utils.RenderServerErrorPage(w)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")
		isDonation := r.FormValue("is_donation") == "on"

		user, err := repository.GetUserByID(sessionUser.ID)
		if err != nil {
			log.Println("CreatePostHandler: Error fetching user:", err)
			utils.RenderServerErrorPage(w)
			return
		}

		var imageFilename string
		file, header, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imageFilename, err = utils.SaveUploadedFile(file, header.Filename, uploadDir)
			if err != nil {
				log.Println("CreatePostHandler: Error saving uploaded image:", err)
				imageFilename = "placeholder.jpg"
			}
		} else {
			log.Println("CreatePostHandler: No image uploaded; using placeholder.")
			imageFilename = "placeholder.jpg"
		}

		if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
			tmpl, _ := template.ParseFiles(postTemplate, navbarTemplate)
			data := viewmodels.CreatePostPageData{
				Error:          "Post title and content cannot be empty or spaces only.",
				Title:          title,
				Content:        content,
				Category:       category,
				ProfilePicture: user.ProfilePicture,
				IsLoggedIn:     true,
			}
			tmpl.Execute(w, data)
			return
		}

		err = repository.CreatePost(user.ID, title, content, category, imageFilename, isDonation, user.Country)
		if err != nil {
			log.Println("CreatePostHandler: Error creating post:", err)
			utils.RenderServerErrorPage(w)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
