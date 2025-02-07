package main

import (
	"database/sql"
	"log"
	"net/http"

	"forum/GlobVar"
	"forum/Handlers"
	middleware "forum/Middleware"
	"forum/Migrations"

	_ "modernc.org/sqlite"
)

func init() {
	var err error
	GlobVar.DB, err = sql.Open("sqlite", "../Database/database.db")
	if err != nil {
		log.Fatal(err)
	}
	Migrations.Migrate()
}
// Upload

func main() {
	defer GlobVar.DB.Close()
	Handlers.HandleUploads()

	// Public routes
	http.Handle("/static/", http.StripPrefix("/static", http.HandlerFunc(Handlers.HandleStatic))) // needs error page
	http.HandleFunc("/", middleware.RateLimiter(Handlers.HandleIndex))
	http.HandleFunc("/post/", middleware.RateLimiter(Handlers.HandlePostPage))
	http.HandleFunc("/Sign_In", middleware.RateLimiter(Handlers.HandleSignIn))
	http.HandleFunc("/Sign_Up", middleware.RateLimiter(Handlers.HandleSignUp))
	http.HandleFunc("/api/auth/status", middleware.RateLimiter(Handlers.HandleAuthStatus))
	http.HandleFunc("/api/checkEmail", middleware.RateLimiter(Handlers.HandleIdentifierDisponibility))
	http.HandleFunc("/api/isValidAuth", middleware.RateLimiter(Handlers.HandleIsValidCredentials))
	

	// Protected routes
	http.HandleFunc("/Comment", middleware.RateLimiter(middleware.ValidateSession(Handlers.HandleComment)))
	http.HandleFunc("/IsLike", middleware.RateLimiter(middleware.ValidateSession(Handlers.HandleLikeDislike)))
	http.HandleFunc("/Log_Out", middleware.RateLimiter(middleware.ValidateSession(Handlers.HandleLogOut)))
	http.HandleFunc("/Profile_Account", middleware.RateLimiter(middleware.ValidateSession(Handlers.HandleProfileAccount)))
	http.HandleFunc("/Update_Profile", middleware.RateLimiter(middleware.ValidateSession(Handlers.HandleProfileUpdate)))
	http.HandleFunc("/New_Post", middleware.RateLimiter(middleware.ValidateSession(Handlers.HandleNewPost)))

	log.Println("server start: http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("server not listener: %v", err)
	}
}