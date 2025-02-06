package Cruds

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"forum/GlobVar"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GenerateUUID generates a new UUID
func GenerateUUID() string {
	id, _ := uuid.NewV4()
	return id.String()
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckUserInfo(email, password string) bool {
	user := GetUser(email)
	if user != nil {
		if !CheckPasswordHash(password, user.PasswordHash) {
			return false
		}	
	} else {
		return false
	}
	return true
}

func GetUserPostCount(userID string) (int, error) {
    query := `
        SELECT COUNT(p.id)
        FROM users u
        LEFT JOIN posts p ON u.id = p.user_id
        WHERE u.id = ?;
    `

    var postCount int
    err := GlobVar.DB.QueryRow(query, userID).Scan(&postCount)
    if err != nil {
        return 0, err
    }

    return postCount, nil
}

func GetUserLikeCount(userID string, isForComment bool) (int, error) {
    query := `
        SELECT COUNT(ld.id)
        FROM users u
        LEFT JOIN posts p ON u.id = p.user_id
        LEFT JOIN likeDislike ld ON p.id = ld.post_id AND ld.is_like = TRUE AND ld.is_comment = ?
        WHERE u.id = ?;
    `

    var likeCount int
    err := GlobVar.DB.QueryRow(query, isForComment, userID).Scan(&likeCount)
    if err != nil {
        return 0, err
    }

    return likeCount, nil
}


func GetUserCommentCount(userID string) (int, error) {
    query := `
        SELECT COUNT(c.id)
        FROM users u
        LEFT JOIN posts p ON u.id = p.user_id
        LEFT JOIN comments c ON p.id = c.post_id
        WHERE u.id = ?;
    `

    var commentCount int
    err := GlobVar.DB.QueryRow(query, userID).Scan(&commentCount)
    if err != nil {
        return 0, err
    }

    return commentCount, nil
}



// Insert Data
func InsertUser(name, image, email, password string) string {
    id := GenerateUUID()
    hashedPassword, err := HashPassword(password)
    if err != nil {
        log.Printf("error hashing password: %v", err)
        return ""
    }
    query := `INSERT INTO users (id, email, user_name, password_hash, user_image) VALUES (?, ?, ?, ?, ?)`
    _, err = GlobVar.DB.Exec(query, id, email, "@"+name, hashedPassword, image)
    if err != nil {
        log.Printf("error inserting user: %v", err)
        return ""
    }
    return id
}

func InsertPost(userId, image, title, content string, categories []string) bool {
    
    // Insert post
    id := GenerateUUID()
    query := `INSERT INTO posts (id, user_id, title, content, image_url) VALUES (?, ?, ?, ?, ?)`
    _, err := GlobVar.DB.Exec(query, id, userId, title, content, image)
    if err != nil {
        log.Printf("error exec query: %v", err)
        return false
    }

    // Insert category if does not exists
    if err = InsertPostCategories(categories, id); err != nil {
        return true    
    }
    return true
}

func InsertPostCategories(categories []string, postId string) error {
    for _, category := range categories {
        query := `SELECT count(*) FROM categories WHERE category_name = ?`
        var exists int
        err := GlobVar.DB.QueryRow(query, category).Scan(&exists)
        if err != nil && err != sql.ErrNoRows {
            fmt.Println(err)
            return err
        }
        if exists == 0 {
            // Insert post
            id := GenerateUUID()
            query := `INSERT INTO categories (id, category_name) VALUES (?, ?)`
            _, err := GlobVar.DB.Exec(query, id, category)
            if err != nil {
                log.Printf("error exec query: %v", err)    
                return err
            }
        }
        query = `INSERT INTO CategoriesByPost (post_id, category_name) VALUES (?, ?)`
        _, err = GlobVar.DB.Exec(query, postId, category)
        if err != nil {
            log.Printf("error exec query: %v", err)
            return err
        }
    }
    return nil
}



func InsertComment(postId, userId, content string) {
	id := GenerateUUID()
	query := `INSERT INTO comments (id, post_id, user_id, content) VALUES (?, ?, ?, ?)`
	_, err := GlobVar.DB.Exec(query, id, postId, userId, content)
	if err != nil {
		log.Printf("error exec query: %v", err)
		return
	}
	log.Println("comment added succesfully")
}

func GetCategories() ([]GlobVar.Categories, error) {
	var categories []GlobVar.Categories

	query := `
		SELECT 
			id, 
			category_name
		FROM 
			categories
	`

	rows, err := GlobVar.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var categorie GlobVar.Categories
		err := rows.Scan(
			&categorie.ID,
			&categorie.CategoryName,
		)
		if err != nil {
			log.Printf("Error scanning categorie row: %v", err)
			continue
		}
		categories = append(categories, categorie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}


func InsertLikeDislike(userId, postId string, isLike bool, isForComment bool) {
	id := GenerateUUID()
	query := `INSERT INTO likeDislike (id, user_id, post_id, is_like, is_comment) VALUES (?, ?, ?, ?, ?)`
	_, err := GlobVar.DB.Exec(query, id, userId, postId, isLike, isForComment)
	if err != nil {
		log.Printf("error exec query: %v", err)
		return
	}
}

func GetPostComments(postId string) ([]GlobVar.Comment, error) {
	var comments []GlobVar.Comment

	query := `
		SELECT 
			comments.id, 
			comments.post_id, 
			comments.user_id, 
			comments.content,
			comments.created_at, 
			comments.updated_at,
			users.user_name AS UserName
		FROM 
			comments
		JOIN 
			users 
		ON 
			comments.user_id = users.id
		WHERE 
			comments.post_id = ?;
	`

	rows, err := GlobVar.DB.Query(query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment GlobVar.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostId,
			&comment.UserId,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.UserName,
		)
		if err != nil {
			log.Printf("Error scanning comment row: %v", err)
			continue
		}
		commentLikes, commentDislikes, err := GetLikesDislikesByPost(comment.ID, true)
		comment.CommentLikes = commentLikes
		comment.CommentDislikes = commentDislikes
		if err != nil {
			log.Printf("Error scanning comment row: %v", err)
			continue
		}
		
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}



func CheckUserLikeDislikeExists(userId, postId string, isForComment bool) (bool, bool) {
	query := `SELECT is_like FROM likeDislike WHERE user_id = ? AND post_id = ? AND is_comment = ?`
	row := GlobVar.DB.QueryRow(query, userId, postId, isForComment)

	var isLike bool
	err := row.Scan(&isLike)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, false
		}
		log.Printf("Error checking like/dislike: %v", err)
		return false, false
	}

	return true, isLike
}


func UpdateLikeDislike(userId, postId string, isLike bool, isForComment bool) {
	query := `UPDATE likeDislike SET is_like = ? WHERE user_id = ? AND post_id = ? AND is_comment = ?`
	_, err := GlobVar.DB.Exec(query, isLike, userId, postId, isForComment)
	if err != nil {
		log.Printf("Error updating like/dislike: %v", err)
	}
}



// this function deletes the like and dislike from the database
func DeleteLikeDislike(userId, postId string, isForComment bool) {
	query := `DELETE FROM likeDislike WHERE user_id = ? AND post_id = ? AND is_comment = ?`
	_, err := GlobVar.DB.Exec(query, userId, postId, isForComment)
	if err != nil {
		log.Printf("error deleting like or dislike: %v", err)
		return
	}
}


// // this function checks if the like or dislike already exist
// func CheckLikeDislikeExists(userId, postId string, isForComment bool) (bool, bool) {
// 	var isLike bool
// 	query := `SELECT is_like FROM likeDislike WHERE user_id = ? AND post_id = ? AND is_comment = ?`
// 	err := GlobVar.DB.QueryRow(query, userId, postId, isForComment).Scan(&isLike)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return false, false
// 		}
// 		return false, false
// 	}
// 	return true, isLike
// }

func GetPostByID(postID string) (*GlobVar.Post, error) {
    query := `SELECT id, user_id, image_url, title, content, created_at FROM posts WHERE id = ?`
    var post GlobVar.Post
    err := GlobVar.DB.QueryRow(query, postID).Scan(&post.ID, &post.UserId, &post.Image, &post.Title, &post.Content, &post.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("post not found")
        }
        return nil, err
    }

	holder, err := GetPostCategoriesByPostId(post.ID)
	post.Categories = holder
	if err != nil {
		return nil, fmt.Errorf("cateogories not found")
	}
    return &post, nil
}

// Update Data
func UpdateUser(email, name, image, password, userId string) error {
	var err error
	var hashedPassword string
	if len(password) != 0 {
		hashedPassword, err = HashPassword(password)
		if err != nil {
			log.Printf("error hashing password: %v", err)
			return err
		}
		query := `UPDATE users SET user_name = ?, user_image = ?, email = ?, password_hash = ? WHERE id = ?`
		_, err = GlobVar.DB.Exec(query, "@"+name, image, email, hashedPassword, userId)
		if err != nil {
			log.Printf("error exec query Update: %v", err)
			return err
		}
	} else {
		query := `UPDATE users SET user_name = ?, user_image = ?, email = ? WHERE id = ?`
		_, err = GlobVar.DB.Exec(query, "@"+name, image, email, userId)
		if err != nil {
			log.Printf("error exec query Update: %v", err)
			return err
		}
	}
	return nil
}

// Get User
func GetUser(value string) *GlobVar.User {
	
	var user GlobVar.User
	query := `SELECT id, email, user_name, password_hash, user_image, created_at FROM users WHERE id = ? OR email = ? OR user_name = ?`

	err := GlobVar.DB.QueryRow(query, value, value, value).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Image, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return nil
	}
	user.UserPostCount, _ = GetUserPostCount(user.ID)
    user.UserCommentCount, _ = GetUserCommentCount(user.ID)
    user.UserLikeCount, _ = GetUserLikeCount(user.ID, false)
    return &user
}



// Get All Data
func GetAllUsers() ([]GlobVar.User, error) {
	query := `SELECT id, email, user_name, password_hash, user_image, created_at FROM users`
	rows, err := GlobVar.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []GlobVar.User
	for rows.Next() {
		var user GlobVar.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Image, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func GetAllPosts() ([]GlobVar.Post, error) {
	query := `SELECT id, user_id, image_url, title, content, created_at FROM posts`
	
	rows, err := GlobVar.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []GlobVar.Post
	for rows.Next() {
		var post GlobVar.Post
		if err := rows.Scan(&post.ID, &post.UserId, &post.Image, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			return nil, err
		}
		
		holder, err := GetPostCategoriesByPostId(post.ID)
		post.Categories = holder

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostCategoriesByPostId(postId string) ([]string, error) {
	categoriesByPost := []string{}
	query := `SELECT category_name FROM CategoriesByPost WHERE post_id = ?`
	rows, err := GlobVar.DB.Query(query, postId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		buffer := ""
		if err := rows.Scan(&buffer); err != nil {
			return nil, err
		}

		categoriesByPost = append(categoriesByPost, buffer)
	}
	
	return categoriesByPost, nil
}

func GetAllComments() ([]GlobVar.Comment, error) {
    query := `
        SELECT c.id, c.post_id, c.user_id, c.content, u.user_name 
        FROM comments c
        JOIN users u ON c.user_id = u.id
    `
    rows, err := GlobVar.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []GlobVar.Comment
    for rows.Next() {
        var cmt GlobVar.Comment
        if err := rows.Scan(&cmt.ID, &cmt.PostId, &cmt.UserId, &cmt.Content, &cmt.UserName); err != nil {
            return nil, err
        }
        comments = append(comments, cmt)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return comments, nil
}

func GetAllLikeDislike() ([]GlobVar.LikeDislike, error) {
	query := `SELECT id, user_id, post_id, is_like FROM likeDislike WHERE is_comment = FALSE`
	rows, err := GlobVar.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likeDislike []GlobVar.LikeDislike
	for rows.Next() {
		var lk GlobVar.LikeDislike
		if err := rows.Scan(&lk.ID, &lk.UserId, &lk.PostId, &lk.IsLike); err != nil {
			return nil, err
		}
		likeDislike = append(likeDislike, lk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return likeDislike, nil
}

func GetLikesDislikesByPost(id string, isForComment bool) (int, int, error) {

	var likes, dislikes int
	queryLikes := `SELECT count(*) FROM likeDislike WHERE post_id = ? and is_like = 1 and is_comment = ?`
	queryDislikes := `SELECT count(*) FROM likeDislike WHERE post_id = ? and is_like = 0 and is_comment = ?`

	err := GlobVar.DB.QueryRow(queryLikes, id, isForComment).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}

	err = GlobVar.DB.QueryRow(queryDislikes, id, isForComment).Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}

func GetCommentsCountByPost(id string) (int, error) {

	var commentsCount int
	query := `SELECT count(*) FROM comments WHERE post_id = ?`
	
	err := GlobVar.DB.QueryRow(query, id).Scan(&commentsCount)
	if err != nil {
		return 0, err
	}

	return commentsCount, nil
}

func ValidateSessionIDAndGetUserID(sessionID string) (string, bool) {
    var expiresAt time.Time
    var userID string
    query := `SELECT user_id, expires_at FROM Session WHERE id = ?`
    err := GlobVar.DB.QueryRow(query, sessionID).Scan(&userID, &expiresAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", false
        }
        log.Printf("Error validating session ID: %v", err)
        return "", false
    }

    // Check if the session is expired
    if time.Now().After(expiresAt) {
        fmt.Println("Session expired:", sessionID)
        return "", false
    }


    return userID, true
}

// ----------- CUSTOM ERROR -----------

type Error struct {
    Status int
    Message string
}

// Function to render error pages with an HTTP status code
func ShowError(w http.ResponseWriter, message string, status int) {

    // Parse the error template
    tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "ErrPage.html"))
    if err != nil {
        // If template parsing fails, fallback to a generic error response
        ShowError(w, "Could not load error page", http.StatusInternalServerError)
        return
    }

	// Set the HTTP status code
    w.WriteHeader(status)
	
    httpError := Error{
        Status: status,
        Message: message,
    }
    // Execute the template with the error message
    err = tmpl.Execute(w, httpError)
    if err != nil {
        // If template execution fails, send a fallback error response with plain text
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal Server Error: Could not render error page"))
    }
}
