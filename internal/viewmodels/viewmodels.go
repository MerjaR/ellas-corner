package viewmodels

import "ellas-corner/internal/repository"

type HomePageData struct {
	IsLoggedIn             bool
	ProfilePicture         string
	ShowConsentBanner      bool
	TopPosts               []repository.Post
	Posts                  []repository.Post
	Categories             []string
	ShowCommentFormForPost int
	ShowEditControls       bool
	ErrorMessage           string
}
