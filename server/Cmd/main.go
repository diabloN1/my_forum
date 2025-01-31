package main

import (
	"database/sql"
	"log"
	"net/http"

	"forum/GlobVar"
	"forum/Handlers"
	middleware "forum/Middleware"
	"forum/Migrations"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	var err error
	GlobVar.DB, err = sql.Open("sqlite3", "../Database/database.db")
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
	http.HandleFunc("/", Handlers.HandleIndex)
	http.HandleFunc("/Sign_In", Handlers.HandleSignIn)
	http.HandleFunc("/Sign_Up", Handlers.HandleSignUp)
	http.HandleFunc("/api/auth/status", Handlers.HandleAuthStatus)

	// Protected routes
	http.HandleFunc("/Comment", middleware.ValidateSession(Handlers.HandleComment))
	http.HandleFunc("/IsLike", middleware.ValidateSession(Handlers.HandleLikeDislike))
	http.HandleFunc("/post/", middleware.ValidateSession(Handlers.HandlePostPage))
	http.HandleFunc("/Log_Out", middleware.ValidateSession(Handlers.HandleLogOut))
	http.HandleFunc("/Profile_Account", middleware.ValidateSession(Handlers.HandleProfileAccount))
	http.HandleFunc("/Update_Profile", middleware.ValidateSession(Handlers.HandleProfileUpdate))
	http.HandleFunc("/New_Post", middleware.ValidateSession(Handlers.HandleNewPost))

	log.Println("server start: http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("server not listener: %v", err)
	}
}