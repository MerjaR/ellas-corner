package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"ellas-corner/internal/viewmodels"
	"html/template"
	"log"
	"net/http"
)

// PostsHandler displays all posts or a single post with its comments
func PostsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("PostsHandler: Request received")

	postID := r.URL.Query().Get("id")

	// Check if the user is logged in
	sessionCookie, err := r.Cookie("session_token")
	isLoggedIn := false
	var userID int
	var currentUser repository.User

	if err == nil {
		userID, err = repository.GetUserIDBySession(sessionCookie.Value)
		if err == nil && userID != 0 {
			isLoggedIn = true
			currentUser, err = repository.GetUserByID(userID)
			if err != nil {
				log.Println("PostsHandler: Failed to fetch user info:", err)
				utils.RenderServerErrorPage(w)
				return
			}
		}
	}

	// Handle single post view
	if postID != "" {
		post, err := repository.GetPostByID(postID, userID)
		if err != nil {
			log.Println("PostsHandler: Error fetching post:", err)
			utils.RenderServerErrorPage(w)
			return
		}

		comments, err := repository.FetchCommentsForPost(post.ID, userID)
		if err != nil {
			log.Printf("PostsHandler: Error fetching comments for post %d: %v", post.ID, err)
			utils.RenderServerErrorPage(w)
			return
		}
		post.Comments = comments

		post.Likes, post.Dislikes, err = repository.FetchReactionsCount(post.ID)
		if err != nil {
			log.Println("PostsHandler: Error fetching reactions count:", err)
			utils.RenderServerErrorPage(w)
			return
		}

		if isLoggedIn {
			userReaction, err := repository.FetchUserReaction(userID, post.ID)
			if err == nil {
				post.UserReaction = userReaction
			}
		}

		tmpl, err := template.ParseFiles("web/templates/post.html")
		if err != nil {
			log.Println("PostsHandler: Error loading template:", err)
			utils.RenderServerErrorPage(w)
			return
		}

		if err := tmpl.Execute(w, post); err != nil {
			log.Println("PostsHandler: Error executing template:", err)
			utils.RenderServerErrorPage(w)
		}
		return
	}

	// Handle all posts view
	posts, err := repository.FetchPosts(userID)
	if err != nil {
		log.Println("PostsHandler: Error fetching posts:", err)
		utils.RenderServerErrorPage(w)
		return
	}

	for i := range posts {
		likes, dislikes, err := repository.FetchReactionsCount(posts[i].ID)
		if err != nil {
			log.Println("PostsHandler: Error fetching reactions count:", err)
			utils.RenderServerErrorPage(w)
			return
		}
		posts[i].Likes = likes
		posts[i].Dislikes = dislikes

		if isLoggedIn {
			userReaction, err := repository.FetchUserReaction(userID, posts[i].ID)
			if err == nil {
				posts[i].UserReaction = userReaction
			}

			if posts[i].IsDonation {
				if currentUser.ShowDonationsInCountryOnly {
					posts[i].ShowDonatedLabel = posts[i].DonationCountry == currentUser.Country &&
						posts[i].DonationCountry != "" &&
						posts[i].DonationCountry != "no_location"
				} else {
					posts[i].ShowDonatedLabel = true
				}
			}
		}
	}

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Println("PostsHandler: Error loading index template:", err)
		utils.RenderServerErrorPage(w)
		return
	}

	data := viewmodels.HomePageData{
		Posts:      posts,
		IsLoggedIn: isLoggedIn,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("PostsHandler: Error executing index template:", err)
		utils.RenderServerErrorPage(w)
	}
}
