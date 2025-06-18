package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HomeHandler: Request received")

	// Only handle root "/"
	if r.URL.Path != "/" {
		// Serve the custom 404 page
		tmpl, err := template.ParseFiles("web/templates/404.html")
		if err != nil {
			log.Println("Error loading 404 template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		// Set the 404 status code
		w.WriteHeader(http.StatusNotFound)

		// Execute the 404 template
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println("Error rendering 404 template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
		}

		return
	}

	// Variables for tracking user session and consent banner
	isLoggedIn := false
	showConsentBanner := true
	var userID int

	// Step 1: Check for session token
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		// If no session token exists, create one for guest users
		sessionToken := utils.GenerateSessionToken()
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/", // Ensure the session token applies site-wide
		})
		log.Println("Generated new session token for guest user")
	} else {
		// Validate session token and check if user is logged in
		userID, err = repository.GetUserIDBySession(sessionCookie.Value)
		if err == nil && userID != 0 {
			isLoggedIn = true
			log.Printf("User is logged in with ID: %d", userID)

			// Check if the user has given cookie consent (optional logic)
			consentGiven, err := repository.CheckCookieConsent(userID)
			if err == nil && consentGiven {
				showConsentBanner = false // Don't show banner if consent is already given
			}
		} else {
			log.Println("Session token is invalid or user ID not found.")
		}
	}

	// Step 2: For non-logged-in users, check if the consent_given cookie is present
	if !isLoggedIn {
		consentCookie, err := r.Cookie("consent_given")
		if err == nil && consentCookie.Value == "true" {
			showConsentBanner = false // Hide banner if consent is already given
		}
	}

	var profilePicture string
	if isLoggedIn {
		user, err := repository.GetUserByID(userID)
		if err != nil {
			log.Println("Error fetching user profile:", err)
		} else {
			profilePicture = user.ProfilePicture
		}
	}

	// Get the post ID from the query to display the comment form
	showCommentFormForPostStr := r.URL.Query().Get("showCommentFormForPost")
	var showCommentFormForPost int
	if showCommentFormForPostStr != "" {
		showCommentFormForPost, _ = strconv.Atoi(showCommentFormForPostStr) // Convert to integer
	}

	// Step 3: Fetch posts from the database (ensure profile pictures are fetched)
	posts, err := repository.FetchPosts() // FetchPosts should now include ProfilePicture for each post
	if err != nil {
		log.Println("Error loading posts:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}
	log.Println("HomeHandler: Posts fetched successfully with profile pictures")

	topPosts, err := repository.FetchTopPostsByLikes(5)
	if err != nil {
		log.Println("Error fetching top liked posts:", err)
		topPosts = []repository.Post{}
	} else {
		for i := range topPosts {
			topPosts[i].Likes, topPosts[i].Dislikes, _ = repository.FetchReactionsCount(topPosts[i].ID)
			if isLoggedIn {
				topPosts[i].UserReaction, _ = repository.FetchUserReaction(userID, topPosts[i].ID)
			}
		}
	}

	// Fetch categories
	categories, err := repository.FetchCategories()
	if err != nil {
		log.Println("Error fetching categories:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}
	log.Println("HomeHandler: Categories fetched:", categories)

	// Step 3.1: Fetch the reaction counts and user reactions (if logged in)
	for i := range posts {
		// Fetch the number of likes and dislikes for each post
		likes, dislikes, err := repository.FetchReactionsCount(posts[i].ID)
		if err != nil {
			log.Println("Error fetching reactions count for post:", posts[i].ID, err)
		}
		posts[i].Likes = likes
		posts[i].Dislikes = dislikes

		// If the user is logged in, fetch the user's reaction to this post
		if isLoggedIn {
			userReaction, err := repository.FetchUserReaction(userID, posts[i].ID)
			if err != nil {
				log.Println("Error fetching user reaction for post:", posts[i].ID, err)
			}
			posts[i].UserReaction = userReaction
		}
	}

	// Step 4: Parse the index template
	tmpl, err := template.ParseFiles("web/templates/index.html", "web/templates/partials/navbar.html", "web/templates/partials/post.html")
	if err != nil {
		log.Println("HomeHandler: Error parsing template", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}
	log.Println("HomeHandler: Template parsed successfully")

	// Step 5: Render the template with the correct data
	data := map[string]interface{}{
		"isLoggedIn":             isLoggedIn,             // Pass the login status to the template
		"ShowConsentBanner":      showConsentBanner,      // Pass consent banner visibility status
		"Posts":                  posts,                  // Pass the fetched posts
		"ShowCommentFormForPost": showCommentFormForPost, // Pass the post ID for displaying the comment form
		"Categories":             categories,
		"ProfilePicture":         profilePicture,
		"TopPosts":               topPosts,
	}
	log.Printf("isLoggedIn: %v\n", isLoggedIn) // Log if the user is logged in

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("HomeHandler: Error executing template", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}
	log.Println("HomeHandler: Template executed successfully")
}
