package handlers

import (
	"ellas-corner/internal/repository"
	"ellas-corner/internal/utils"
	"html/template"
	"log"
	"net/http"
)

// PostsHandler fetches all posts or a single post with its comments
func PostsHandler(w http.ResponseWriter, r *http.Request) {
	// Check if a post ID is provided (for a single post view)
	postID := r.URL.Query().Get("id")

	// Check if the user is logged in
	sessionCookie, err := r.Cookie("session_token")
	isLoggedIn := false
	var userID int

	if err == nil {
		// Validate session token and check if user is logged in
		userID, err = repository.GetUserIDBySession(sessionCookie.Value)
		if err == nil && userID != 0 {
			isLoggedIn = true
		}
	}

	var currentUser repository.User
	if isLoggedIn {
		currentUser, err = repository.GetUserByID(userID)
		if err != nil {
			log.Println("Failed to fetch user info for donation filtering:", err)
		}
	}

	if postID != "" {
		// Fetch the specific post by ID, including its comments (handled in GetPostByID)
		post, err := repository.GetPostByID(postID)
		if err != nil {
			log.Println("Error fetching post:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		// Fetch comments for this post and attach them to the post structure
		comments, err := repository.FetchCommentsForPost(post.ID)
		if err != nil {
			log.Printf("Error fetching comments for post %d: %v", post.ID, err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		// Attach comments to the post
		post.Comments = comments

		// Fetch like/dislike counts for this post
		likes, dislikes, err := repository.FetchReactionsCount(post.ID)
		if err != nil {
			log.Println("Error fetching post reactions:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		// Attach the like and dislike counts
		post.Likes = likes
		post.Dislikes = dislikes

		// If user is logged in, check if they have already reacted
		if isLoggedIn {
			userReaction, err := repository.FetchUserReaction(userID, post.ID)
			if err == nil {
				post.UserReaction = userReaction
			}
		}

		// Render the single post with comments and reactions
		tmpl, err := template.ParseFiles("web/templates/post.html")
		if err != nil {
			log.Println("Could not load template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		// Execute the template and pass the post data (with comments, reactions) to the template
		err = tmpl.Execute(w, post)
		if err != nil {
			log.Println("Error executing template:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
		}
		return
	}

	// Fetch all posts if no post ID is provided
	posts, err := repository.FetchPosts()
	if err != nil {
		log.Println("Error fetching posts:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	// Add like/dislike counts and user reactions (if logged in) for each post
	for i := range posts {
		likes, dislikes, err := repository.FetchReactionsCount(posts[i].ID)
		if err != nil {
			log.Println("Error fetching post reactions:", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderServerErrorPage(w)
			return
		}

		// Attach like and dislike counts to each post
		posts[i].Likes = likes
		posts[i].Dislikes = dislikes

		// If user is logged in, check if they have already reacted to each post
		if isLoggedIn {
			userReaction, err := repository.FetchUserReaction(userID, posts[i].ID)
			if err == nil {
				posts[i].UserReaction = userReaction
			}
		}

		if isLoggedIn && posts[i].IsDonation {
			if currentUser.ShowDonationsInCountryOnly {
				// ✅ Show donation tag only if it has a valid country AND matches user’s country
				if posts[i].DonationCountry != "" &&
					posts[i].DonationCountry != "no_location" &&
					posts[i].DonationCountry == currentUser.Country {
					posts[i].ShowDonatedLabel = true
				} else {
					posts[i].ShowDonatedLabel = false
				}
			} else {
				// ✅ If user wants to see all donations regardless of country
				posts[i].ShowDonatedLabel = true
			}
		} else {
			// ✅ Not a donation or user not logged in
			posts[i].ShowDonatedLabel = false
		}

		log.Printf("Post %d | IsDonation: %v | Country: %s | Show: %v\n",
			posts[i].ID,
			posts[i].IsDonation,
			posts[i].DonationCountry,
			posts[i].ShowDonatedLabel)
	}

	// Render the posts in the index template with like/dislike functionality
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Println("Could not load template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
		return
	}

	// Execute the template and pass the posts data (with comments, reactions) to the template
	err = tmpl.Execute(w, map[string]interface{}{
		"Posts":      posts, // Using the unified Post structure with comments, reactions
		"IsLoggedIn": isLoggedIn,
	})
	if err != nil {
		log.Println("Error executing template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.RenderServerErrorPage(w)
	}
}
