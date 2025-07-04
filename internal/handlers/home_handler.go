package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"ellas-corner/internal/viewmodels"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HomeHandler: Request received")

	// Handle only root path
	if r.URL.Path != "/" {
		tmpl, err := template.ParseFiles("web/templates/404.html")
		if err != nil {
			log.Println("HomeHandler: Error loading 404 template:", err)
			utils.RenderServerErrorPage(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		if err := tmpl.Execute(w, nil); err != nil {
			log.Println("HomeHandler: Error rendering 404 template:", err)
			utils.RenderServerErrorPage(w)
		}
		return
	}

	// Initialize session/user state
	isLoggedIn := false
	showConsentBanner := true
	var userID int
	var currentUser repository.User
	var profilePicture string

	// Step 1: Check for existing session token or create guest session
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		sessionToken := utils.GenerateSessionToken()
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})
		log.Println("HomeHandler: Generated session token for guest")
	} else {
		userID, err = repository.GetUserIDBySession(sessionCookie.Value)
		if err == nil && userID != 0 {
			isLoggedIn = true
			log.Printf("HomeHandler: Logged-in user ID: %d", userID)

			consentGiven, err := repository.CheckCookieConsent(userID)
			if err == nil && consentGiven {
				showConsentBanner = false
			}
		} else {
			log.Println("HomeHandler: Invalid session token or user not found")
		}
	}

	// Step 2: Fallback for cookie-based consent (non-logged-in users)
	if !isLoggedIn {
		if consentCookie, err := r.Cookie("consent_given"); err == nil && consentCookie.Value == "true" {
			showConsentBanner = false
		}
	}

	// Step 3: Load current user (if logged in)
	if isLoggedIn {
		user, err := repository.GetUserByID(userID)
		if err != nil {
			log.Println("HomeHandler: Error fetching user profile:", err)
		} else {
			currentUser = user
			profilePicture = user.ProfilePicture
		}
	}

	// Get query param to reveal comment form
	showCommentFormForPostStr := r.URL.Query().Get("showCommentFormForPost")
	showCommentFormForPost, _ := strconv.Atoi(showCommentFormForPostStr)

	// Step 4: Fetch all posts
	posts, err := repository.FetchPosts(userID)
	if err != nil {
		log.Println("HomeHandler: Error fetching posts:", err)
		utils.RenderServerErrorPage(w)
		return
	}
	log.Println("HomeHandler: Posts fetched")

	// Step 5: Fetch top liked posts
	topPosts, err := repository.FetchTopPostsByLikes(5)
	if err != nil {
		log.Println("HomeHandler: Error fetching top liked posts:", err)
		topPosts = []repository.Post{}
	} else {
		for i := range topPosts {
			topPosts[i].Likes, topPosts[i].Dislikes, _ = repository.FetchReactionsCount(topPosts[i].ID)
			if isLoggedIn {
				topPosts[i].UserReaction, _ = repository.FetchUserReaction(userID, topPosts[i].ID)
			}
		}
	}

	// Step 6: Fetch categories
	categories, err := repository.FetchCategories()
	if err != nil {
		log.Println("HomeHandler: Error fetching categories:", err)
		utils.RenderServerErrorPage(w)
		return
	}
	log.Println("HomeHandler: Categories fetched")

	// Step 7: Populate reaction and donation visibility for each post
	for i := range posts {
		likes, dislikes, err := repository.FetchReactionsCount(posts[i].ID)
		if err != nil {
			log.Println("HomeHandler: Error fetching reactions for post", posts[i].ID, ":", err)
		}
		posts[i].Likes = likes
		posts[i].Dislikes = dislikes

		if isLoggedIn {
			userReaction, err := repository.FetchUserReaction(userID, posts[i].ID)
			if err != nil {
				log.Println("HomeHandler: Error fetching user reaction for post", posts[i].ID, ":", err)
			}
			posts[i].UserReaction = userReaction

			// Donation visibility filtering for logged-in users
			if posts[i].IsDonation {
				if currentUser.ShowDonationsInCountryOnly {
					posts[i].ShowDonatedLabel = (posts[i].DonationCountry == currentUser.Country)
				} else {
					posts[i].ShowDonatedLabel = true
				}
			}
		} else {
			// For guests, show donation label if post is a donation
			if posts[i].IsDonation {
				posts[i].ShowDonatedLabel = true
			}
		}
	}

	// Step 8: Render the homepage
	tmpl, err := template.ParseFiles(
		"web/templates/index.html",
		"web/templates/partials/navbar.html",
		"web/templates/partials/post.html",
	)
	if err != nil {
		log.Println("HomeHandler: Error parsing template:", err)
		utils.RenderServerErrorPage(w)
		return
	}
	log.Println("HomeHandler: Templates parsed successfully")

	data := viewmodels.HomePageData{
		IsLoggedIn:             isLoggedIn,
		ProfilePicture:         profilePicture,
		ShowConsentBanner:      showConsentBanner,
		TopPosts:               topPosts,
		Posts:                  posts,
		Categories:             categories,
		ShowCommentFormForPost: showCommentFormForPost,
		ShowEditControls:       false, // Future logic for author-owned posts
		ErrorMessage:           "",    // Optional error display
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("HomeHandler: Error executing template:", err)
		utils.RenderServerErrorPage(w)
		return
	}
	log.Println("HomeHandler: Page rendered successfully")
}
