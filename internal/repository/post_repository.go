package repository

import (
	"database/sql"
	"ellas-corner/internal/db"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Post struct {
	ID                 int
	UserID             int
	Username           string
	ProfilePicture     string // New field for the user's profile picture
	Title              string
	Content            string
	Category           string
	CreatedAt          string // Original CreatedAt as string
	FormattedCreatedAt string // Formatted CreatedAt for display
	Comments           []Comment
	Likes              int
	Dislikes           int
	UserReaction       string
	Image              string
	ShowDonatedLabel   bool
	IsDonation         bool
	DonationCountry    string
}

// CreatePost saves a new post to the database
func CreatePost(userID int, title, content, category, image string, isDonation bool, donationCountry string) error {
	query := `
	INSERT INTO posts (user_id, title, content, category, image, is_donation, donation_country)
	VALUES (?, ?, ?, ?, ?, ?, ?)
`
	_, err := db.DB.Exec(query, userID, title, content, category, image, isDonation, donationCountry)
	if err != nil {
		log.Println("Error creating post:", err)
		return err
	}
	return nil
}

func FetchPosts() ([]Post, error) {
	query := `
        SELECT posts.id, posts.title, posts.content, posts.user_id, posts.category, posts.created_at, users.username, users.profile_picture, COALESCE(posts.image, '') AS image, posts.is_donation, COALESCE(posts.donation_country, '') AS donation_country

        FROM posts
        JOIN users ON posts.user_id = users.id
        ORDER BY posts.created_at DESC`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Println("Error fetching posts:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Category, &post.CreatedAt, &post.Username, &post.ProfilePicture, &post.Image, &post.IsDonation, &post.DonationCountry)
		if err != nil {
			log.Println("Error scanning post:", err)
			return nil, err
		}

		// Parse and format CreatedAt
		parsedTime, err := time.Parse(time.RFC3339, post.CreatedAt)
		if err != nil {
			log.Println("Error parsing CreatedAt:", err)
			post.FormattedCreatedAt = post.CreatedAt // Fallback if parsing fails
		} else {
			post.FormattedCreatedAt = parsedTime.Format("02 Jan 2006, 15:04")
		}

		// Fetch comments
		comments, err := FetchCommentsForPost(post.ID)
		if err != nil {
			log.Println("Error fetching comments for post:", err)
			return nil, err
		}
		post.Comments = comments
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostByID(postID string) (*Post, error) {
	query := `
        SELECT posts.id, posts.title, posts.content, posts.user_id, posts.category, posts.created_at, users.username, users.profile_picture, COALESCE(posts.image, '') AS image

        FROM posts
        JOIN users ON posts.user_id = users.id
        WHERE posts.id = ?`

	var post Post
	err := db.DB.QueryRow(query, postID).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.Category,
		&post.CreatedAt,
		&post.Username,
		&post.ProfilePicture,
		&post.Image)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No post found with ID: %s", postID)
			return nil, nil
		}
		log.Println("Error fetching post by ID:", err)
		return nil, err
	}

	// Parse and format CreatedAt
	parsedTime, err := time.Parse(time.RFC3339, post.CreatedAt)
	if err != nil {
		log.Println("Error parsing CreatedAt:", err)
		post.FormattedCreatedAt = post.CreatedAt // Fallback if parsing fails
	} else {
		post.FormattedCreatedAt = parsedTime.Format("02 Jan 2006, 15:04")
	}

	// Fetch comments
	comments, err := FetchCommentsForPost(post.ID)
	if err != nil {
		log.Println("Error fetching comments for post:", err)
		return nil, err
	}
	post.Comments = comments

	return &post, nil
}

// AddReaction adds or updates a user's reaction to a post
func AddReaction(userID int, postID int, reactionType string) error {
	// First, check if the user already reacted to this post
	query := `SELECT reaction_type FROM post_reactions WHERE post_id = ? AND user_id = ?`
	var existingReaction string
	err := db.DB.QueryRow(query, postID, userID).Scan(&existingReaction)

	if err == sql.ErrNoRows {
		// No previous reaction, insert a new one
		insertQuery := `INSERT INTO post_reactions (post_id, user_id, reaction_type) VALUES (?, ?, ?)`
		log.Printf("AddReaction: Inserting new reaction for userID=%d, postID=%d, reactionType=%s", userID, postID, reactionType)
		_, err := db.DB.Exec(insertQuery, postID, userID, reactionType)
		return err
	} else if err != nil {
		log.Println("AddReaction: Error checking reaction:", err)
		return err
	}

	// If the user has already reacted, update the reaction
	log.Printf("AddReaction: Updating reaction for userID=%d, postID=%d, existingReaction=%s, newReaction=%s", userID, postID, existingReaction, reactionType)
	if existingReaction != reactionType {
		updateQuery := `UPDATE post_reactions SET reaction_type = ? WHERE post_id = ? AND user_id = ?`
		_, err = db.DB.Exec(updateQuery, reactionType, postID, userID)
		return err
	}

	// If the user has already reacted with the same type, no action needed
	log.Printf("AddReaction: No update needed, userID=%d, postID=%d, reactionType=%s", userID, postID, reactionType)
	return nil
}

// FetchReactionsCount fetches the count of likes and dislikes for a specific post
func FetchReactionsCount(postID int) (likes int, dislikes int, err error) {
	query := `
        SELECT 
            COALESCE(SUM(CASE WHEN reaction_type = 'like' THEN 1 ELSE 0 END), 0) AS likes,
            COALESCE(SUM(CASE WHEN reaction_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes
        FROM post_reactions WHERE post_id = ?`
	err = db.DB.QueryRow(query, postID).Scan(&likes, &dislikes)
	if err != nil {
		log.Println("Error fetching reactions count:", err)
	}
	return likes, dislikes, err
}

// FetchUserReaction retrieves the reaction ("like" or "dislike") for a specific post by a specific user
func FetchUserReaction(userID int, postID int) (string, error) {
	var reaction string

	// SQL query to check if the user has reacted to the post
	query := `SELECT reaction_type FROM post_reactions WHERE user_id = ? AND post_id = ?`

	err := db.DB.QueryRow(query, userID, postID).Scan(&reaction)
	if err == sql.ErrNoRows {
		// No reaction found, return an empty string
		return "", nil
	} else if err != nil {
		log.Println("Error fetching user reaction:", err)
		return "", err
	}

	return reaction, nil
}

//Filtering

func FetchCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM posts ORDER BY category`
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Println("Error fetching categories from DB:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			log.Println("Error scanning category row:", err)
			return nil, err
		}
		log.Println("Fetched category:", category) // Log each fetched category
		categories = append(categories, category)
	}

	// Check if any categories were fetched
	if len(categories) == 0 {
		log.Println("No categories found in database.")
	}

	return categories, nil
}

// FetchFilteredPosts retrieves posts based on the given filter criteria.
func FetchFilteredPosts(category, createdPosts, likedPosts, startDate, endDate string, userID int, isLoggedIn bool) ([]Post, error) {
	query := `
		SELECT posts.id, posts.title, posts.content, posts.user_id, posts.category, posts.created_at,
		       users.username, users.profile_picture, COALESCE(posts.image, '') AS image
		FROM posts
		JOIN users ON posts.user_id = users.id
		WHERE 1=1`

	// Dynamic filters
	if category != "" {
		query += " AND posts.category = '" + category + "'"
	}
	if startDate != "" {
		query += " AND DATE(posts.created_at) >= DATE('" + startDate + "')"
	}
	if endDate != "" {
		query += " AND DATE(posts.created_at) <= DATE('" + endDate + "')"
	}
	if createdPosts == "true" && isLoggedIn {
		query += " AND posts.user_id = " + strconv.Itoa(userID)
	}
	if likedPosts == "true" && isLoggedIn {
		query += ` AND posts.id IN (SELECT post_id FROM post_reactions WHERE user_id = ` + strconv.Itoa(userID) + ` AND reaction_type = 'like')`
	}

	query += " ORDER BY posts.created_at DESC"

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Println("Error fetching filtered posts:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var createdAt time.Time

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Category, &createdAt, &post.Username, &post.ProfilePicture, &post.Image)
		if err != nil {
			log.Println("Error scanning filtered post:", err)
			return nil, err
		}

		post.CreatedAt = createdAt.Format(time.RFC3339)
		post.FormattedCreatedAt = createdAt.Format("02 Jan 2006, 15:04")

		// ✅ Fetch comments
		comments, err := FetchCommentsForPost(post.ID)
		if err != nil {
			log.Println("Error fetching comments:", err)
			return nil, err
		}
		post.Comments = comments

		// ✅ Fetch reactions
		likes, dislikes, err := FetchReactionsCount(post.ID)
		if err != nil {
			log.Println("Error fetching reactions count:", err)
			return nil, err
		}
		post.Likes = likes
		post.Dislikes = dislikes

		// ✅ If logged in, fetch user reaction
		if isLoggedIn {
			reaction, err := FetchUserReaction(userID, post.ID)
			if err != nil {
				log.Println("Error fetching user reaction:", err)
				return nil, err
			}
			post.UserReaction = reaction
		}

		posts = append(posts, post)
	}

	return posts, nil
}

//Search

func SearchPosts(searchQuery string, userID int, isLoggedIn bool) ([]Post, error) {
	query := `
    SELECT posts.id, posts.title, posts.content, posts.category, posts.created_at,
           users.username, users.profile_picture, COALESCE(posts.image, '') AS image
    FROM posts
    JOIN users ON posts.user_id = users.id
    WHERE posts.title LIKE '%' || ? || '%' 
       OR posts.content LIKE '%' || ? || '%'
       OR posts.category LIKE '%' || ? || '%'
       OR users.username LIKE '%' || ? || '%'
    ORDER BY posts.created_at DESC`

	rows, err := db.DB.Query(query, searchQuery, searchQuery, searchQuery, searchQuery)
	if err != nil {
		log.Println("Error searching posts:", err)
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
			&post.Image,
		)
		if err != nil {
			log.Println("Error scanning post during search:", err)
			return nil, err
		}

		// Format date
		parsedTime, err := time.Parse(time.RFC3339, post.CreatedAt)
		if err != nil {
			log.Println("Error parsing CreatedAt:", err)
			post.FormattedCreatedAt = post.CreatedAt
		} else {
			post.FormattedCreatedAt = parsedTime.Format("02 Jan 2006, 15:04")
		}

		// Fetch comments
		comments, err := FetchCommentsForPost(post.ID)
		if err != nil {
			log.Println("Error fetching comments for post:", err)
			return nil, err
		}
		post.Comments = comments

		// Fetch reaction counts
		likes, dislikes, err := FetchReactionsCount(post.ID)
		if err != nil {
			log.Println("Error fetching reaction counts:", err)
			return nil, err
		}
		post.Likes = likes
		post.Dislikes = dislikes

		// If logged in, get user's reaction
		if isLoggedIn {
			userReaction, err := FetchUserReaction(userID, post.ID)
			if err != nil {
				log.Println("Error fetching user reaction for post:", err)
				return nil, err
			}
			post.UserReaction = userReaction
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// DeletePost removes a post from the database by its ID
func DeletePost(postID int) error {
	query := "DELETE FROM posts WHERE id = ?"

	_, err := db.DB.Exec(query, postID)
	if err != nil {
		log.Println("Error deleting post:", err)
		return err
	}

	return nil
}

// UpdatePost updates a post's title, content, and category in the database
func UpdatePost(postID int, title, content, category string, isDonation bool) error {
	query := `
		UPDATE posts
		SET title = ?, content = ?, category = ?, is_donation = ?
		WHERE id = ?
	`
	_, err := db.DB.Exec(query, title, content, category, isDonation, postID)
	return err
}

func FetchLikedPosts(userID int) ([]Post, error) {
	query := `
		SELECT posts.id, posts.title, posts.content, posts.user_id, posts.category, posts.created_at,
		       users.username, users.profile_picture
		FROM posts
		JOIN post_reactions ON posts.id = post_reactions.post_id
		JOIN users ON posts.user_id = users.id
		WHERE post_reactions.user_id = ? AND post_reactions.reaction_type = 'like'
		ORDER BY posts.created_at DESC
	`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		log.Println("Error fetching liked posts:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var createdAt time.Time

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.UserID,
			&post.Category,
			&createdAt,
			&post.Username,
			&post.ProfilePicture,
		)
		if err != nil {
			log.Println("Error scanning liked post:", err)
			return nil, err
		}

		post.CreatedAt = createdAt.Format(time.RFC3339)
		post.FormattedCreatedAt = createdAt.Format("02 Jan 2006, 15:04")

		posts = append(posts, post)
	}

	return posts, nil
}

// UpdatePostWithImage updates a post's title, content, category, donation flag, and image in the database
func UpdatePostWithImage(postID int, title, content, category string, isDonation bool, image string) error {
	query := `
		UPDATE posts
		SET title = ?, content = ?, category = ?, is_donation = ?, image = ?
		WHERE id = ?
	`
	_, err := db.DB.Exec(query, title, content, category, isDonation, image, postID)
	if err != nil {
		log.Println("Error updating post with image:", err)
	}
	return err
}

func FetchTopPostsByLikes(limit int) ([]Post, error) {
	query := `
		SELECT posts.id, posts.title, posts.content, posts.category, posts.image, posts.created_at,
		       users.username, users.profile_picture
		FROM posts
		JOIN users ON posts.user_id = users.id
		LEFT JOIN post_reactions AS likes ON posts.id = likes.post_id AND likes.reaction_type = 'like'
		GROUP BY posts.id
		ORDER BY COUNT(likes.id) DESC
		LIMIT ?`

	rows, err := db.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var createdAt time.Time
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Image, &createdAt, &post.Username, &post.ProfilePicture)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = createdAt.Format(time.RFC3339)
		post.FormattedCreatedAt = createdAt.Format("02 Jan 2006, 15:04")
		posts = append(posts, post)
	}
	return posts, nil
}

func SaveImageFile(file multipart.File, handler *multipart.FileHeader) (string, error) {
	// Ensure the target directory exists
	err := os.MkdirAll("web/static/uploads", os.ModePerm)
	if err != nil {
		return "", err
	}

	// Use the original filename
	filename := filepath.Base(handler.Filename)
	dstPath := filepath.Join("web/static/uploads", filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return filename, nil
}
