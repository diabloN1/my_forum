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
	cookies "forum/cookies"

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

func GetUserLikeCount(userID string) (int, error) {
    query := `
        SELECT COUNT(ld.id)
        FROM users u
        LEFT JOIN posts p ON u.id = p.user_id
        LEFT JOIN likeDislike ld ON p.id = ld.post_id AND ld.is_like = TRUE
        WHERE u.id = ?;
    `

    var likeCount int
    err := GlobVar.DB.QueryRow(query, userID).Scan(&likeCount)
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
    _, err = GlobVar.DB.Exec(query, id, email, name, hashedPassword, image)
    if err != nil {
        log.Printf("error inserting user: %v", err)
        return ""
    }
    return id
}

func InsertPost(userId, image, title, content, category string) bool {
	id := GenerateUUID()
	query := `INSERT INTO posts (id, user_id, title, content, image_url, category) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := GlobVar.DB.Exec(query, id, userId, title, content, image, category)
	if err != nil {
		log.Printf("error exec query: %v", err)
		return false
	}
	return true
}

func InsertComment(postId, userId, content string) {
	id := GenerateUUID()
	query := `INSERT INTO comments (id, post_id, user_id, content) VALUES (?, ?, ?, ?)`
	_, err := GlobVar.DB.Exec(query, id, postId, userId, content)
	if err != nil {
		log.Printf("error exec query: %v", err)
		return
	}
}

func InsertCategory(byUserId, categoryName string) {
	id := GenerateUUID()
	query := `INSERT INTO categories (id, category_name, created_by_user_id) VALUES (?, ?, ?)`
	_, err := GlobVar.DB.Exec(query, id, categoryName, byUserId)
	if err != nil {
		log.Printf("error exec query: %v", err)
		return
	}
}

func InsertLikeDislike(userId, postId string, isLike bool) {
	id := GenerateUUID()
	query := `INSERT INTO likeDislike (id, user_id, post_id, is_like) VALUES (?, ?, ?, ?)`
	_, err := GlobVar.DB.Exec(query, id, userId, postId, isLike)
	if err != nil {
		log.Printf("error exec query: %v", err)
		return
	}
}


func CheckUserLikeDislikeExists(userId, postId string) (bool, bool) {
	query := `SELECT is_like FROM likeDislike WHERE user_id = ? AND post_id = ?`
	row := GlobVar.DB.QueryRow(query, userId, postId)

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


func UpdateLikeDislike(userId, postId string, isLike bool) {
	query := `UPDATE likeDislike SET is_like = ? WHERE user_id = ? AND post_id = ?`
	_, err := GlobVar.DB.Exec(query, isLike, userId, postId)
	if err != nil {
		log.Printf("Error updating like/dislike: %v", err)
	}
}



// this function deletes the like and dislike from the database
func DeleteLikeDislike(userId, postId string) {
	query := `DELETE FROM likeDislike WHERE user_id = ? AND post_id = ?`
	_, err := GlobVar.DB.Exec(query, userId, postId)
	if err != nil {
		log.Printf("error deleting like or dislike: %v", err)
		return
	}
}


// this function checks if the like or dislike already exist
func CheckLikeDislikeExists(userId, postId string) (bool, bool) {
	var isLike bool
	query := `SELECT is_like FROM likeDislike WHERE user_id = ? AND post_id = ?`
	err := GlobVar.DB.QueryRow(query, userId, postId).Scan(&isLike)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, false
		}
		return false, false
	}
	return true, isLike
}

func GetPostByID(postID string) (string, *GlobVar.Post, error) {
    query := `SELECT id, user_id, image_url, title, content, category, created_at FROM posts WHERE id = ?`
    var post GlobVar.Post
    err := GlobVar.DB.QueryRow(query, postID).Scan(&post.ID, &post.UserId, &post.Image, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return "" ,nil, fmt.Errorf("post not found")
        }
        return "" ,nil, err
    }
    return "" ,&post, nil
}

// Update Data
func UpdateUser(email, name, image, password, userId string) {
	var err error
	var hashedPassword string
	if len(password) != 0 {
		hashedPassword, err = HashPassword(password)
		if err != nil {
			log.Printf("error hashing password: %v", err)
			return
		}
		query := `UPDATE users SET user_name = ?, user_image = ?, email = ?, password_hash = ? WHERE id = ?`
		_, err = GlobVar.DB.Exec(query, name, image, email, hashedPassword, userId)
		if err != nil {
			log.Printf("error exec query Update: %v", err)
		}
	} else {
		query := `UPDATE users SET user_name = ?, user_image = ?, email = ? WHERE id = ?`
		_, err = GlobVar.DB.Exec(query, name, image, email, userId)
		if err != nil {
			log.Printf("error exec query Update: %v", err)
		}
	}
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
    user.UserLikeCount, _ = GetUserLikeCount(user.ID)
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
	query := `SELECT id, user_id, image_url, title, content, category, created_at FROM posts`
	
	rows, err := GlobVar.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []GlobVar.Post
	for rows.Next() {
		var post GlobVar.Post
		if err := rows.Scan(&post.ID, &post.UserId, &post.Image, &post.Title, &post.Content, &post.Category, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
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
	query := `SELECT id, user_id, post_id, is_like FROM likeDislike`
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

func GetLikesDislikesByPost(id string) (int, int, error) {

	var likes, dislikes int
	queryLikes := `SELECT count(*) FROM likeDislike WHERE post_id = ? and is_like = 1`
	queryDislikes := `SELECT count(*) FROM likeDislike WHERE post_id = ? and is_like = 0`

	err := GlobVar.DB.QueryRow(queryLikes, id).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}

	err = GlobVar.DB.QueryRow(queryDislikes, id).Scan(&dislikes)
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

func Set_Cookies_Handler(w http.ResponseWriter, r *http.Request, userID string) {
    var sessionID, token string
    var err error

    // Generate a unique session ID
    sessionID, err = cookies.Generate_Cookie_session()
    if err != nil {
        ShowError(w, "Internal Server Error", http.StatusInternalServerError)
        log.Printf("Error generating session ID: %v", err)
        return
    }

    // Generate a unique token for the session
    for {
        token, err = cookies.Generate_Cookie_session()
        if err != nil {
            ShowError(w, "Internal Server Error", http.StatusInternalServerError)
            log.Printf("Error generating session token: %v", err)
            return
        }

        // Check if the token already exists in the database
        var exists bool
        query := `SELECT EXISTS(SELECT 1 FROM Session WHERE token = ?)`
        err = GlobVar.DB.QueryRow(query, token).Scan(&exists)
        if err != nil {
            ShowError(w, "Internal Server Error", http.StatusInternalServerError)
            log.Printf("Error checking token existence: %v", err)
            return
        }

        if !exists {
            break // Token is unique, exit the loop
        }
    }

    // Insert the session into the database
    expiresAt := time.Now().Add(7 * 24 * time.Hour) // Session expires in 7 days
    query := `INSERT INTO Session (id, user_id, token, expires_at) VALUES (?, ?, ?, ?)`
    _, err = GlobVar.DB.Exec(query, sessionID, userID, token, expiresAt)
    if err != nil {
        log.Printf("Error storing session in database: %v", err) // Debugging
        ShowError(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    fmt.Println("Session created for user:", userID) // Debugging

    // Set the session cookie
    cookie := &http.Cookie{
        Name:     "Session_ID",
        Value:    sessionID,
        Path:     "/",
        Secure:   true,
        HttpOnly: true,
        Expires:  expiresAt,
        SameSite: http.SameSiteStrictMode,
    }
    http.SetCookie(w, cookie)
}

func Delete_Cookie_Handler(w http.ResponseWriter, r *http.Request) {
	// Get the session cookie
	cookie, err := r.Cookie("Session_ID")
	if err != nil {
		// No session cookie found
		http.Redirect(w, r, "/Sign_In", http.StatusSeeOther)
		return
	}

	// Delete the session from the database
	query := `DELETE FROM Session WHERE id = ?`
	_, err = GlobVar.DB.Exec(query, cookie.Value)
	if err != nil {
		ShowError(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error deleting session from database: %v", err)
		return
	}

	// Clear the session cookie
	cookie = &http.Cookie{
		Name:     "Session_ID",
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour), // Expire the cookie
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}


// ----------- CUSTOM ERROR -----------

type Error struct {
    Status int
    Message string
}

// Function to render error pages with an HTTP status code
func ShowError(w http.ResponseWriter, message string, status int) {

    // Set the HTTP status code
    w.WriteHeader(status)

    // Parse the error template
    tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "ErrPage.html"))
    if err != nil {
        // If template parsing fails, fallback to a generic error response
        ShowError(w, "Could not load error page", http.StatusInternalServerError)
        return
    }

    httpError := Error{
        Status: status,
        Message: message,
    }
    // Execute the template with the error message
    err = tmpl.Execute(w, httpError)
    if err != nil {
        // If template execution fails, respond with a generic error
        ShowError(w, "Could not render error page", http.StatusInternalServerError)
    }
}






// package cruds

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	"forum/GlobVar"

// 	"github.com/gofrs/uuid"
// 	"golang.org/x/crypto/bcrypt"
// )

// // Utility Functions

// // GenerateUUID generates a new UUID.
// func GenerateUUID() string {
// 	id, _ := uuid.NewV4()
// 	return id.String()
// }

// // HashPassword hashes a password using bcrypt.
// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return "", fmt.Errorf("error hashing password: %w", err)
// 	}
// 	return string(bytes), nil
// }

// // ExecQuery executes a query and logs errors.
// func ExecQuery(query string, args ...interface{}) error {
// 	_, err := GlobVar.DB.Exec(query, args...)
// 	if err != nil {
// 		log.Printf("Error executing query (%s): %v", query, err)
// 		return err
// 	}
// 	return nil
// }

// // QueryRow scans a single row and logs errors.
// func QueryRow(query string, dest ...interface{}) error {
// 	err := GlobVar.DB.QueryRow(query).Scan(dest...)
// 	if err != nil && err != sql.ErrNoRows {
// 		log.Printf("Error querying row (%s): %v", query, err)
// 	}
// 	return err
// }

// // CRUD Operations

// // InsertUser inserts a new user into the database.
// func InsertUser(name, image, email, password string) (string, error) {
// 	id := GenerateUUID()
// 	hashedPassword, err := HashPassword(password)
// 	if err != nil {
// 		return "", err
// 	}
// 	query := `INSERT INTO users (id, email, user_name, password_hash, user_image) VALUES (?, ?, ?, ?, ?)`
// 	return id, ExecQuery(query, id, email, name, hashedPassword, image)
// }

// // InsertPost inserts a new post into the database.
// func InsertPost(userID, image, title, content, category string) error {
// 	id := GenerateUUID()
// 	query := `INSERT INTO posts (id, user_id, title, content, image_url, category) VALUES (?, ?, ?, ?, ?, ?)`
// 	return ExecQuery(query, id, userID, title, content, image, category)
// }

// // InsertComment inserts a new comment into the database.
// func InsertComment(postID, userID, content string) error {
// 	id := GenerateUUID()
// 	query := `INSERT INTO comments (id, post_id, user_id, content) VALUES (?, ?, ?, ?)`
// 	return ExecQuery(query, id, postID, userID, content)
// }

// // InsertCategory inserts a new category into the database.
// func InsertCategory(userID, categoryName string) error {
// 	id := GenerateUUID()
// 	query := `INSERT INTO categories (id, category_name, created_by_user_id) VALUES (?, ?, ?)`
// 	return ExecQuery(query, id, categoryName, userID)
// }

// // InsertLikeDislike inserts a like or dislike into the database.
// func InsertLikeDislike(userID, postID string, isLike bool) error {
// 	id := GenerateUUID()
// 	query := `INSERT INTO likeDislike (id, user_id, post_id, is_like) VALUES (?, ?, ?, ?)`
// 	return ExecQuery(query, id, userID, postID, isLike)
// }

// // DeleteLikeDislike deletes a like or dislike from the database.
// func DeleteLikeDislike(userID, postID string) error {
// 	query := `DELETE FROM likeDislike WHERE user_id = ? AND post_id = ?`
// 	return ExecQuery(query, userID, postID)
// }

// // CheckLikeDislikeExists checks if a like or dislike exists.
// func CheckLikeDislikeExists(userID, postID string) (bool, bool, error) {
// 	var isLike bool
// 	query := `SELECT is_like FROM likeDislike WHERE user_id = ? AND post_id = ?`
// 	err := GlobVar.DB.QueryRow(query, userID, postID).Scan(&isLike)
// 	if err == sql.ErrNoRows {
// 		return false, false, nil
// 	} else if err != nil {
// 		return false, false, err
// 	}
// 	return true, isLike, nil
// }

// // GetPostByID retrieves a post by its ID.
// func GetPostByID(postID string) (*GlobVar.Post, error) {
// 	query := `SELECT id, user_id, image_url, title, content, category, created_at FROM posts WHERE id = ?`
// 	var post GlobVar.Post
// 	err := QueryRow(query, &post.ID, &post.UserId, &post.Image, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
// 	if err == sql.ErrNoRows {
// 		return nil, fmt.Errorf("post not found")
// 	}
// 	return &post, err
// }

// // UpdateUser updates a user's information.
// func UpdateUser(email, name, image, password, userID string) error {
// 	query := `UPDATE users SET user_name = ?, user_image = ?, email = ?, password_hash = ? WHERE id = ?`
// 	hashedPassword, err := HashPassword(password)
// 	if err != nil {
// 		return err
// 	}
// 	return ExecQuery(query, name, image, email, hashedPassword, userID)
// }

// // Bulk Data Retrieval

// // GetAllUsers retrieves all users.
// func GetAllUsers() ([]GlobVar.User, error) {
// 	query := `SELECT id, email, user_name, password_hash, user_image, created_at FROM users`
// 	return fetchUsers(query)
// }

// // fetchUsers scans multiple users from a query result.
// func fetchUsers(query string) ([]GlobVar.User, error) {
// 	rows, err := GlobVar.DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var users []GlobVar.User
// 	for rows.Next() {
// 		var user GlobVar.User
// 		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Image, &user.CreatedAt); err != nil {
// 			return nil, err
// 		}
// 		users = append(users, user)
// 	}
// 	return users, rows.Err()
// }
