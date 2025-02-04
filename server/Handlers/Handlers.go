package Handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	Cruds "forum/Api"
	"forum/GlobVar"
	cookies "forum/cookies"
)

func HandleStatic(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
        Cruds.ShowError(w, "404- Not Found", 404)
        return
    }

    fs := http.FileServer(http.Dir(GlobVar.StaticPath))
	fs.ServeHTTP(w, r)
}

func HandleUploads() {
	http.Handle("/Uploads/", http.StripPrefix("/Uploads/", http.FileServer(http.Dir("../Uploads"))))
}

func HandlePostPage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/post/" {
        Cruds.ShowError(w, "404 not found", http.StatusNotFound)
        return
    }

    // Extract the post ID from the URL
    postID := r.URL.Query().Get("id")
    if postID == "" {
        Cruds.ShowError(w, "Post ID is required", http.StatusBadRequest)
        return
    }

    // Fetch the post details
    post, err := Cruds.GetPostByID(postID)

    if err != nil {
        Cruds.ShowError(w, "Post not found", http.StatusNotFound)
        return
    }

    // Fetch the user who created the post
    user := Cruds.GetUser(post.UserId)
    if user == nil {
        Cruds.ShowError(w, "User not found", http.StatusNotFound)
        return
    }

    // Fetch comments for the post
    postComments, err := Cruds.GetPostComments(postID)
    if err != nil {
        Cruds.ShowError(w, "Failed to fetch comments", http.StatusInternalServerError)
        return
    }

    // Fetch likes and dislikes for the post
    likes, dislikes, err := Cruds.GetLikesDislikesByPost(postID, false)
    if err != nil {
        Cruds.ShowError(w, "Failed to fetch likes/dislikes", http.StatusInternalServerError)
    }

    // Prepare the data to be passed to the template
    data := struct {
        Post     *GlobVar.Post
        User     *GlobVar.User
        Comments []GlobVar.Comment
        Likes    int
        Dislikes int
    }{
        Post:     post,
        User:     user,
        Comments: postComments,
        Likes:    likes,
        Dislikes: dislikes,
    }

    // Render the post page
    tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "post_page.html"))
    if err != nil {
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }
}

func HandleComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Comment" {
		Cruds.ShowError(w, "404", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
        Cruds.ShowError(w, "405", http.StatusMethodNotAllowed)
		return
	}

    comment := r.FormValue("content")
    postId := r.FormValue("postId")
    userId := r.Context().Value("userID").(string)
    if strings.TrimSpace(comment) == "" || len(comment) > 2000 {
        http.Redirect(w, r, "/post/?id="+postId, http.StatusSeeOther)
        return 
    }
    Cruds.InsertComment(postId, userId, comment)
    http.Redirect(w, r, "/post/?id="+postId, http.StatusSeeOther)

}

func HandleLikeDislike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/IsLike" {
		Cruds.ShowError(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		// Extracting values from the form
		postId := r.FormValue("postId")
		commentId := r.FormValue("commentId")
		userId, _ := r.Context().Value("userID").(string)
		isLike := r.FormValue("isLike") == "true"
        isForComment := r.FormValue("isComment") == "true"
        postToRedirect := postId


        if isForComment {
            postId = commentId
        }

		// Validate inputs
		if (postId == "") || userId == "" {
			Cruds.ShowError(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Check if the user already liked/disliked this post
		exists, currentIsLike := Cruds.CheckUserLikeDislikeExists(userId, postId, isForComment)

		if exists {
			if isLike == currentIsLike {
				// If the current action matches the existing action, remove the like/dislike
				Cruds.DeleteLikeDislike(userId, postId, isForComment)
			} else {
				// If the current action is different, update the like/dislike
				Cruds.UpdateLikeDislike(userId, postId, isLike, isForComment)
			}
		} else {
			// If no record exists, insert a new like/dislike
			Cruds.InsertLikeDislike(userId, postId, isLike, isForComment)
		}

		// Redirect the user back to the previous page
		referer := r.Referer()
		if referer == "http://localhost:8080/" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/post/?id=" + postToRedirect, http.StatusSeeOther)
		}
		return
	}

	// Handle unsupported methods
	Cruds.ShowError(w, "404 - Page Not Found", http.StatusNotFound)
}


func HandleLogOut(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Log_Out" {
		Cruds.ShowError(w, "404", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		// Delete the session cookie and session from the database
		Delete_Cookie_Handler(w, r)

		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	Cruds.ShowError(w, "404 - Page Not Found", http.StatusNotFound)
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/Sign_In" {
        Cruds.ShowError(w, "404 not found", http.StatusNotFound)
        return
    }

    if r.Method == http.MethodPost {

        email := r.FormValue("email")
        password := r.FormValue("password")
        

        if email == "" || len(password) < 8 || len(email) > 50 || len(password) > 20 {
            Cruds.ShowError(w, "500", http.StatusInternalServerError)
        } 

        // Fetch the user
        user := Cruds.GetUser(email)

        if user == nil {
            http.Redirect(w, r, "/Sign_In", http.StatusSeeOther)
            return
        }

        // Compare the password
        if !Cruds.CheckPasswordHash(password, user.PasswordHash) {
            http.Redirect(w, r, "/Sign_In", http.StatusSeeOther)
            return
        }

        // Set the session cookie
        Set_Cookies_Handler(w, r, user.ID)

        // Redirect to home page
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    
    // Render the sign-in page
    tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "sign-in-page.html"))
    if err != nil {
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, nil)
    if err != nil {
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }
}

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/Sign_Up" {
        Cruds.ShowError(w, "404 - Page Not Found", http.StatusNotFound)
        return
    }

    if r.Method == http.MethodPost {
        name := r.FormValue("name")
        email := r.FormValue("email")
        password := r.FormValue("password")
        passwordConfirmation := r.FormValue("passwordConfirmation")

        image := GlobVar.DefaultImage
        
        // Check email and name availability
        u1 := Cruds.GetUser(email)
        u2 := Cruds.GetUser(name)

        emailRegxp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
        isValidEmail := u1 == nil && emailRegxp.MatchString(email) && email != "" && len(password) < 20
        isValidName := u2 == nil && !strings.Contains(name, "@") && !strings.Contains(name, " ") && name != "" && len(name) < 20
        isValidPassword := len(password) < 8 || password != passwordConfirmation || len(password) > 20 || len(passwordConfirmation) > 20
        if !isValidName || !isValidEmail || isValidPassword {
            http.Redirect(w, r, "/Sign_Up", http.StatusSeeOther)       
            return
        }

        // Insert the new user into the database
        userID := Cruds.InsertUser(name, image, email, password)
        if userID == "" {
            http.Redirect(w, r, "/Sign_Up", http.StatusSeeOther)
            return
        }

        // Set the session cookie for the new user
        Set_Cookies_Handler(w, r, userID)
        http.Redirect(w, r, "/", http.StatusSeeOther)
		 
        return
    }

    if r.Method != http.MethodGet {
        Cruds.ShowError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    // Render the sign-up page
    tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "sign-up-page.html"))
    if err != nil {
        log.Printf("Error parsing template: %v", err)
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, nil)
    if err != nil {
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }

}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Cruds.ShowError(w, "404", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		Cruds.ShowError(w, "405", http.StatusMethodNotAllowed)
		return
	}

	posts, err := Cruds.GetAllPosts()
	if err != nil {
		Cruds.ShowError(w, "500 -11", http.StatusInternalServerError)
		return
	}
    fmt.Println(posts)
     
	if len(posts) > 0 {
		for i := range posts {
			//user
            user := Cruds.GetUser(posts[i].UserId)
            posts[i].UserName = user.Name
            posts[i].UserImage = user.Image

			//comment
			posts[i].NbrComment, err = Cruds.GetCommentsCountByPost(posts[i].ID)
            if err != nil {
                Cruds.ShowError(w, "500", http.StatusBadRequest)
                return
            }

			//likedislike
            posts[i].NbrLike, posts[i].NbrDislike, err = Cruds.GetLikesDislikesByPost(posts[i].ID, false)
            if err != nil {
                Cruds.ShowError(w, "500", http.StatusBadRequest)
                return
            }
		}
	}

	tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "index.html"))
	if err != nil {
		Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
    
    fmt.Println(posts)
	err = tmpl.Execute(w, posts)
	if err != nil {
		Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
	}
}

func HandleProfileAccount(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/Profile_Account" {
        Cruds.ShowError(w, "404 - Page Not Found", http.StatusNotFound)
        return
    }

    if r.Method != http.MethodGet {
        Cruds.ShowError(w, "405 - Method Not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Retrieve the user ID from the context
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        Cruds.ShowError(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Query the user using the user ID from the context
    data := Cruds.GetUser(userID)
    

    if data == nil {
        Cruds.ShowError(w, "User not found", http.StatusNotFound)
        return
    }

    tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "account-page.html"))
    if err != nil {
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }
}

func HandleProfileUpdate(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/Update_Profile" {
        Cruds.ShowError(w, "page - not found", 404)
        return
    }

    // Retrieve the user ID from the context
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        Cruds.ShowError(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
	

    data := Cruds.GetUser(userID)

    if r.Method == http.MethodPost {
        name := r.FormValue("name")
        email := r.FormValue("email")
        password := r.FormValue("password")
        passwordConfirmation := r.FormValue("passwordConfirmation")
    	imagePath := data.Image // Default to existing image

        err := r.ParseMultipartForm(20 * 1024 * 1024) // 20mb
        // To be Impelented !!!!!!!!!!!!!
        if err != nil {
            Cruds.ShowError(w, "size image abouve than 20MB",http.StatusBadRequest)
            return
        }
        // Handle file upload
        file, fileHeader, _ := r.FormFile("image")

        if file != nil {
            defer file.Close()
            copyFile, err := os.Create("../Uploads/" + fileHeader.Filename)
            if err != nil {
                Cruds.ShowError(w, "err open file", http.StatusInternalServerError)
                return
            }
            defer copyFile.Close()
            hold := make([]byte, fileHeader.Size)
            
            _, err = file.Read(hold)
            if err != nil {
                http.Error(w, "err copy file to newFile", http.StatusInternalServerError)
                return
            }

            _,err = copyFile.Write(hold)
            if err != nil {
                http.Error(w, "err copy file to newFile", http.StatusInternalServerError)
                return
            }
            imagePath = "/Uploads/" + fileHeader.Filename
        }
        // Check email and name availability
        u1 := Cruds.GetUser(email)
        u2 := Cruds.GetUser(name)

        if len(name) == 0 {
            name = data.Name
            u2 = nil
        }
        if len(email) == 0 {
            email = data.Email
            u1 = nil
        }
		if len(password) == 0 || len(name) > 20 {
			password = ""
		}
        
        emailRegxp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
        isValidEmail := u1 == nil && emailRegxp.MatchString(email) && len(password) < 20
        isValidName := u2 == nil && !strings.Contains(name, "@") && !strings.Contains(name, " ") && len(name) < 20
        isValidPassword := password == "" || password == passwordConfirmation
        if (!isValidEmail || !isValidName || !isValidPassword) {
            http.Redirect(w, r, "/Update_Profile", http.StatusSeeOther)       
            return
        }
        
		// Update user in the database
		err = Cruds.UpdateUser(email, name, imagePath, password, userID)
        if err != nil {
            Cruds.ShowError(w, "500", 500)
        }
		http.Redirect(w, r, "/Profile_Account", http.StatusSeeOther)
		return		
    }

    tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "update-account-page.html"))
    if err != nil {
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    data.Name = data.Name[1:]
    err = tmpl.Execute(w, data)
    if err != nil {
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }
}

func HandleNewPost(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/New_Post" {
        Cruds.ShowError(w, "404 not found", http.StatusNotFound)
        return
    }
	userID, ok := r.Context().Value("userID").(string)
    if !ok {
        Cruds.ShowError(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    if r.Method == http.MethodPost {
        data := Cruds.GetUser(userID)
        title := r.FormValue("title")
        categories := strings.Split(r.FormValue("categories"), " ")
        content := r.FormValue("content")

        // [start] upload image
        image := GlobVar.DefaultImage
        err := r.ParseMultipartForm(20 * 1024 * 1024) // 20 MB
        if err != nil {
            Cruds.ShowError(w, "size image abouve than 20MB",http.StatusBadRequest)
            return
        }
        
        file, fileHeader, err := r.FormFile("post_image")
        if err != nil && err != http.ErrMissingFile {
            http.Error(w, "err formFile",http.StatusBadRequest)
            return
        }
        if file != nil {
            defer file.Close()
            copyFile, err := os.Create("../Uploads/" + fileHeader.Filename)
            if err != nil {
                http.Error(w, "err open file", http.StatusInternalServerError)
                return
            }
            defer copyFile.Close()
            hold := make([]byte, fileHeader.Size)
            file.Read(hold)
            _,err = copyFile.Write(hold)
            if err != nil {
                http.Error(w, "err copy file to newFile", http.StatusInternalServerError)
                return
            }
            image = "/Uploads/"+fileHeader.Filename
        }
        


        isValidInputs := title != "" && len(categories) != 0 && content != "" && len(title) < 50 && len(categories) < 10 && len(content) < 1200
        if isValidInputs && Cruds.InsertPost(data.ID, image, title, content, categories) {
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        } else {
            http.Redirect(w, r, "/New_Post", http.StatusSeeOther)
        }
        return
    }

    tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "new-post-page.html"))
    if err != nil {
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    data, err := Cruds.GetCategories()
    if err != nil {
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        Cruds.ShowError(w, "Internal server error", http.StatusInternalServerError)
        return
    }
}

func Set_Cookies_Handler(w http.ResponseWriter, r *http.Request, userID string) {
	sessionID, err := cookies.Generate_Cookie_session()
	if err != nil {
		Cruds.ShowError(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error generating session ID: %v", err)
		return
	}

	// Insert the session into the database
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // Session expires in 7 days
	query := `INSERT INTO Session (id, user_id, expires_at) VALUES (?, ?, ?)`
	_, err = GlobVar.DB.Exec(query, sessionID, userID, expiresAt)
	if err != nil {
		Cruds.ShowError(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error storing session in database: %v", err)
		return
	}

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

// Delete_Cookie_Handler deletes the session cookie.
func Delete_Cookie_Handler(w http.ResponseWriter, r *http.Request) {
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
		Cruds.ShowError(w, "Internal Server Error", http.StatusInternalServerError)
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

func HandleAuthStatus(w http.ResponseWriter, r *http.Request) {
    isAuthenticated := false
    cookie, err := r.Cookie("Session_ID")
    if err == nil {
        _, isAuthenticated = Cruds.ValidateSessionIDAndGetUserID(cookie.Value)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{
        "isAuthenticated": isAuthenticated,
    })
}

func HandleIdentifierDisponibility(w http.ResponseWriter, r *http.Request) {
    
    identifier := r.FormValue("identifier")

    fmt.Println(identifier)
    user := Cruds.GetUser(identifier)
    isDisponible := user == nil

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{
        "isDisponible": isDisponible,
    })
}

func HandleIsValidCredentials(w http.ResponseWriter, r *http.Request) {
    buffer := make([]byte, r.ContentLength)
    nb, err := r.Body.Read(buffer)
    fmt.Println(nb, err, string(buffer))
    if nb == 0 || err != io.EOF {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]bool{
            "isValid": false,
        })
        return
    }

    var data = make(map[string]string)
    err = json.Unmarshal(buffer, &data)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]bool{
            "isValid": false,
        })
        return
    }
    fmt.Println("data :", data)
    isValid := Cruds.CheckUserInfo(data["email"],data["password"])

    fmt.Println(isValid)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{
        "isValid": isValid,
    })

    // user := Cruds.GetUser(email)
    // isValid := user == nil
    // w.Header().Set("Content-Type", "application/json")
    // json.NewEncoder(w).Encode(map[string]bool{
    //     "isValid": isValid,
    // })
}