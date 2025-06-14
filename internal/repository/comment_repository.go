package repository

import (
	"database/sql"
	"ellas-corner/internal/db"
	"log"
	"strconv"
	"time"
)

// Comment represents a comment on a post
type Comment struct {
	ID                 int
	PostID             int
	UserID             int
	Username           string
	ProfilePicture     string // New field for the user's profile picture
	Content            string
	CreatedAt          string
	FormattedCreatedAt string // Human-readable formatted timestamp
	PostTitle          string
	Likes              int    // Count of likes for the comment
	Dislikes           int    // Count of dislikes for the comment
	UserReaction       string // User's reaction to the comment ("like" or "dislike")
	ParentCommentID    *int   // Nullable field to store parent comment ID
}

// FetchCommentsForPost retrieves comments for a specific post, including the user's profile picture.
func FetchCommentsForPost(postID int) ([]Comment, error) {
	query := `
		SELECT comments.id, comments.post_id, comments.user_id, comments.content, comments.created_at, users.username, users.profile_picture
		FROM comments
		JOIN users ON comments.user_id = users.id
		WHERE comments.post_id = ?
		ORDER BY comments.created_at ASC`

	rows, err := db.DB.Query(query, postID)
	if err != nil {
		log.Println("Error fetching comments for post:", err)
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Username, &comment.ProfilePicture); err != nil {
			log.Println("Error scanning comment:", err)
			return nil, err
		}

		// Parse and format the CreatedAt string
		parsedTime, err := time.Parse(time.RFC3339, comment.CreatedAt)
		if err != nil {
			log.Println("Error parsing CreatedAt for comment:", comment.ID, ":", err)
			comment.FormattedCreatedAt = comment.CreatedAt // Use raw string if parsing fails
		} else {
			comment.FormattedCreatedAt = parsedTime.Format("02 Jan 2006, 15:04") // Desired format
		}

		// Fetch likes and dislikes for the comment
		likes, dislikes, err := FetchCommentReactionsCount(comment.ID)
		if err == nil {
			comment.Likes = likes
			comment.Dislikes = dislikes
		}

		comments = append(comments, comment)
	}
	return comments, nil
}

// CreateComment adds a new comment to the database
func CreateComment(userID int, postIDStr string, content string) error {
	// Convert postIDStr to an integer
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("Error converting post ID to integer:", err)
		return err
	}

	// Insert the comment into the database
	query := "INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)"
	_, err = db.DB.Exec(query, userID, postID, content)
	if err != nil {
		log.Println("Error creating comment:", err)
		return err
	}
	return nil
}

// AddCommentReaction adds or updates a user's reaction to a comment
func AddCommentReaction(userID int, commentID int, reactionType string) error {
	// First, check if the user already reacted to this comment
	query := `SELECT reaction_type FROM comment_reactions WHERE comment_id = ? AND user_id = ?`
	var existingReaction string
	err := db.DB.QueryRow(query, commentID, userID).Scan(&existingReaction)

	if err == sql.ErrNoRows {
		// No previous reaction, insert a new one
		insertQuery := `INSERT INTO comment_reactions (comment_id, user_id, reaction_type) VALUES (?, ?, ?)`
		_, err := db.DB.Exec(insertQuery, commentID, userID, reactionType)
		return err
	} else if err != nil {
		return err
	}

	// If the user has already reacted, update the reaction
	if existingReaction != reactionType {
		updateQuery := `UPDATE comment_reactions SET reaction_type = ? WHERE comment_id = ? AND user_id = ?`
		_, err = db.DB.Exec(updateQuery, reactionType, commentID, userID)
		return err
	}

	// If the user has already reacted with the same type, no action is needed
	return nil
}

// FetchCommentReactionsCount fetches the count of likes and dislikes for a specific comment
func FetchCommentReactionsCount(commentID int) (likes int, dislikes int, err error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN reaction_type = 'like' THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN reaction_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes
		FROM comment_reactions WHERE comment_id = ?`
	err = db.DB.QueryRow(query, commentID).Scan(&likes, &dislikes)
	return likes, dislikes, err
}

// FetchUserCommentReaction retrieves the reaction ("like" or "dislike") for a specific comment by a specific user
func FetchUserCommentReaction(userID int, commentID int) (string, error) {
	var reaction string

	// SQL query to check if the user has reacted to the comment
	query := `SELECT reaction_type FROM comment_reactions WHERE user_id = ? AND comment_id = ?`

	err := db.DB.QueryRow(query, userID, commentID).Scan(&reaction)
	if err == sql.ErrNoRows {
		// No reaction found, return an empty string
		return "", nil
	} else if err != nil {
		return "", err
	}

	return reaction, nil
}

func DeleteComment(commentID int) error {
	query := "DELETE FROM comments WHERE id = ?"
	_, err := db.DB.Exec(query, commentID)
	if err != nil {
		log.Println("Error deleting comment:", err)
	}
	return err
}
