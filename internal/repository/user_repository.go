package repository

import (
	"database/sql"
	"log"
	"time"
)

// User represents the user entity in the database
type User struct {
	ID                         int
	Username                   string
	Email                      string
	Password                   string
	ProfilePicture             string
	Country                    string
	ShowDonationsInCountryOnly bool
	IsDonation                 bool
}

// CreateUser inserts a new user into the database, including the profile picture
func CreateUser(username, email, password, profilePicture string) error {
	query := "INSERT INTO users (username, email, password, profile_picture) VALUES (?, ?, ?, ?)"
	_, err := database.Conn.Exec(query, username, email, password, profilePicture)
	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}
	return nil
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(email string) (*User, error) {
	var user User
	query := "SELECT id, username, email, password, profile_picture FROM users WHERE email = ?"
	var profilePicture sql.NullString // Use sql.NullString to handle NULL values
	err := database.Conn.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &profilePicture)

	// If profile_picture is NULL, assign an empty string or a default picture
	if profilePicture.Valid {
		user.ProfilePicture = profilePicture.String
	} else {
		user.ProfilePicture = "" // Or set a default picture like "/static/default.png"
	}
	if err == sql.ErrNoRows {
		return nil, nil // No user found with that email
	} else if err != nil {
		log.Println("Error fetching user by email:", err)
		return nil, err
	}
	return &user, nil
}

// GetUserIDBySession retrieves the user ID associated with the given session token
func GetUserIDBySession(sessionToken string) (int, error) {
	var userID int
	query := "SELECT user_id FROM sessions WHERE session_token = ?"
	err := database.Conn.QueryRow(query, sessionToken).Scan(&userID)
	if err == sql.ErrNoRows {
		log.Println("No session found for the given token")
		return 0, nil
	} else if err != nil {
		log.Println("Error retrieving user ID by session token:", err)
		return 0, err
	}
	return userID, nil
}

// SaveSessionToken stores the session token in the database
func SaveSessionToken(userID int, sessionToken string) error {
	query := "INSERT INTO sessions (user_id, session_token) VALUES (?, ?) ON CONFLICT(user_id) DO UPDATE SET session_token = ?"
	_, err := database.Conn.Exec(query, userID, sessionToken, sessionToken)
	if err != nil {
		log.Println("Error saving session token:", err)
		return err
	}
	return nil
}

// CheckCookieConsent checks if the user has given consent for cookies
func CheckCookieConsent(userID int) (bool, error) {
	var consentGiven bool
	query := "SELECT consent_given FROM cookie_consent WHERE user_id = ?"
	err := database.Conn.QueryRow(query, userID).Scan(&consentGiven)
	if err == sql.ErrNoRows {
		// No consent record found, treat this as consent not given
		return false, nil
	} else if err != nil {
		// Log any other errors
		log.Println("Error checking cookie consent:", err)
		return false, err
	}
	return consentGiven, nil
}

// SaveCookieConsent saves the user's consent decision
func SaveCookieConsent(userID int, consentGiven bool) error {
	query := "INSERT INTO cookie_consent (user_id, consent_given) VALUES (?, ?) ON CONFLICT(user_id) DO UPDATE SET consent_given = ?"
	_, err := database.Conn.Exec(query, userID, consentGiven, consentGiven)
	if err != nil {
		log.Println("Error saving cookie consent:", err)
		return err
	}
	return nil
}

// DeleteSession removes a session from the database based on the session token
func DeleteSession(sessionToken string) error {
	query := "DELETE FROM sessions WHERE session_token = ?"
	_, err := database.Conn.Exec(query, sessionToken)
	if err != nil {
		log.Println("Error deleting session:", err)
		return err
	}
	return nil
}

//Fetch data for user profile with these functions

// GetUserByID retrieves a user by their ID and handles NULL values for profile_picture
func GetUserByID(userID int) (User, error) {
	query := "SELECT id, username, email, password, profile_picture, country, show_donations_in_country_only FROM users WHERE id = ?"

	var user User
	var profilePicture sql.NullString
	var country sql.NullString
	var showDonations bool

	err := database.Conn.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &profilePicture, &country, &showDonations)
	if err != nil {
		return User{}, err
	}

	// If profile_picture is NULL, assign an empty string or a default picture
	if profilePicture.Valid {
		user.ProfilePicture = profilePicture.String
	} else {
		user.ProfilePicture = "" // Or set a default picture like "/static/default.png"
	}

	if country.Valid {
		user.Country = country.String
	} else {
		user.Country = "no_location"
	}

	user.ShowDonationsInCountryOnly = showDonations

	return user, nil
}

// FetchPostsByUser fetches all posts by a specific user
func FetchPostsByUser(userID int) ([]Post, error) {
	query := `
		SELECT posts.id, posts.title, posts.content, posts.category, posts.created_at,
		       users.username, users.profile_picture, posts.is_donation,
		       COALESCE(posts.image, '') AS image,
		       (SELECT COUNT(*) FROM post_reactions WHERE post_id = posts.id AND reaction_type = 'like') AS likes,
		       (SELECT COUNT(*) FROM post_reactions WHERE post_id = posts.id AND reaction_type = 'dislike') AS dislikes
		FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE posts.user_id = ?
		ORDER BY posts.created_at DESC`

	rows, err := database.Conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Category,
			&post.CreatedAt,
			&post.Username,
			&post.ProfilePicture,
			&post.IsDonation,
			&post.Image,
			&post.Likes,
			&post.Dislikes,
		)
		if err != nil {
			return nil, err
		}

		// Format CreatedAt
		parsedTime, err := time.Parse(time.RFC3339, post.CreatedAt)
		if err == nil {
			post.FormattedCreatedAt = parsedTime.Format("02 Jan 2006, 15:04")
		} else {
			post.FormattedCreatedAt = post.CreatedAt
		}

		// Fetch and attach comments
		comments, err := FetchCommentsForPost(post.ID, userID)

		if err != nil {
			return nil, err
		}
		post.Comments = comments

		// Note: UserReaction is optional here, since you're fetching the user's own posts
		// You can skip it unless you want to display your own likes/dislikes

		posts = append(posts, post)
	}

	return posts, nil
}

// FetchCommentsByUser retrieves all comments made by a specific user, along with the post titles
func FetchCommentsByUser(userID int) ([]Comment, error) {
	query := `
        SELECT comments.id, comments.content, comments.created_at, posts.title,
               users.username, users.profile_picture
        FROM comments 
        JOIN posts ON comments.post_id = posts.id 
        JOIN users ON comments.user_id = users.id
        WHERE comments.user_id = ? 
        ORDER BY comments.created_at DESC`

	rows, err := database.Conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		var createdAt string
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&createdAt,
			&comment.PostTitle,
			&comment.Username,
			&comment.ProfilePicture,
		)
		if err != nil {
			return nil, err
		}

		parsedTime, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			comment.FormattedCreatedAt = createdAt
		} else {
			comment.FormattedCreatedAt = parsedTime.Format("02 Jan 2006, 15:04")
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func FetchLikedPostsByUser(userID int) ([]Post, error) {
	query := `
        SELECT posts.id, posts.title, posts.content, posts.category, posts.created_at,
               users.username, users.profile_picture, COALESCE(posts.image, '') AS image,
               (SELECT COUNT(*) FROM post_reactions WHERE post_id = posts.id AND reaction_type = 'like') AS likes,
               (SELECT COUNT(*) FROM post_reactions WHERE post_id = posts.id AND reaction_type = 'dislike') AS dislikes
        FROM posts
        JOIN post_reactions ON posts.id = post_reactions.post_id
        JOIN users ON posts.user_id = users.id
        WHERE post_reactions.user_id = ? AND post_reactions.reaction_type = 'like'
        ORDER BY posts.created_at DESC`

	rows, err := database.Conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likedPosts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Category,
			&post.CreatedAt,
			&post.Username,
			&post.ProfilePicture,
			&post.Image,
			&post.Likes,
			&post.Dislikes,
		)
		if err != nil {
			return nil, err
		}

		parsedTime, err := time.Parse(time.RFC3339, post.CreatedAt)
		if err == nil {
			post.FormattedCreatedAt = parsedTime.Format("02 Jan 2006, 15:04")
		} else {
			post.FormattedCreatedAt = post.CreatedAt
		}

		comments, err := FetchCommentsForPost(post.ID, userID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		post.UserReaction = "like" // For consistency in the template
		likedPosts = append(likedPosts, post)
	}

	return likedPosts, nil
}

func FetchDislikedPostsByUser(userID int) ([]Post, error) {
	query := `
        SELECT posts.id, posts.title, posts.content, posts.category, posts.created_at,
               users.username, users.profile_picture, COALESCE(posts.image, '') AS image,
               (SELECT COUNT(*) FROM post_reactions WHERE post_id = posts.id AND reaction_type = 'like') AS likes,
               (SELECT COUNT(*) FROM post_reactions WHERE post_id = posts.id AND reaction_type = 'dislike') AS dislikes
        FROM posts
        JOIN post_reactions ON posts.id = post_reactions.post_id
        JOIN users ON posts.user_id = users.id
        WHERE post_reactions.user_id = ? AND post_reactions.reaction_type = 'dislike'
        ORDER BY posts.created_at DESC`

	rows, err := database.Conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dislikedPosts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Category,
			&post.CreatedAt,
			&post.Username,
			&post.ProfilePicture,
			&post.Image,
			&post.Likes,
			&post.Dislikes,
		)
		if err != nil {
			return nil, err
		}

		parsedTime, err := time.Parse(time.RFC3339, post.CreatedAt)
		if err == nil {
			post.FormattedCreatedAt = parsedTime.Format("02 Jan 2006, 15:04")
		} else {
			post.FormattedCreatedAt = post.CreatedAt
		}

		comments, err := FetchCommentsForPost(post.ID, userID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		post.UserReaction = "dislike"
		dislikedPosts = append(dislikedPosts, post)
	}

	return dislikedPosts, nil
}

// UpdateProfilePicture updates the profile picture path in the database
func UpdateProfilePicture(userID int, filePath string) error {
	query := "UPDATE users SET profile_picture = ? WHERE id = ?"
	_, err := database.Conn.Exec(query, filePath, userID)
	if err != nil {
		log.Println("Error updating profile picture:", err)
		return err
	}
	return nil
}

func UpdateUserPreferences(userID int, country string, showDonations bool) error {
	query := `
		UPDATE users 
		SET country = ?, show_donations_in_country_only = ? 
		WHERE id = ?
	`

	_, err := database.Conn.Exec(query, country, showDonations, userID)
	if err != nil {
		log.Println("Error updating user preferences:", err)
		return err
	}

	return nil
}

func UpdateDonationCountryForUser(userID int, newCountry string) error {
	query := `
		UPDATE posts
		SET donation_country = ?
		WHERE user_id = ? AND is_donation = 1
	`
	_, err := database.Conn.Exec(query, newCountry, userID)
	return err
}
