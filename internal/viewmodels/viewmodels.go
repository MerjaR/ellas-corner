package viewmodels

import (
	"ellas-corner/internal/db"
	"ellas-corner/internal/repository"
)

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

type CreatePostPageData struct {
	IsLoggedIn     bool
	ProfilePicture string
	Error          string
	Title          string
	Content        string
	Category       string
}

type EditPostPageData struct {
	IsLoggedIn     bool
	ProfilePicture string
	Post           repository.Post
	Categories     []string
}

type FilterPageData struct {
	IsLoggedIn             bool
	ProfilePicture         string
	Posts                  []repository.Post
	Categories             []string
	Category               string // selected category
	ShowCommentFormForPost int
	ShowEditControls       bool
}

type LikedPostsPageData struct {
	IsLoggedIn     bool
	ProfilePicture string
	Username       string
	LikedPosts     []repository.Post
	CuratedItems   []db.BabyBoxItem
}

type ProfilePageData struct {
	Username                   string
	Email                      string
	ProfilePicture             string
	Country                    string
	ShowDonationsInCountryOnly bool
	IsLoggedIn                 bool
	Posts                      []repository.Post
	Comments                   []repository.Comment
	LikedPosts                 []repository.Post
	DislikedPosts              []repository.Post
}

type SearchPageData struct {
	IsLoggedIn             bool
	ProfilePicture         string
	SearchQuery            string
	Posts                  []repository.Post
	ShowEditControls       bool
	ShowCommentFormForPost int
}
