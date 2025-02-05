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
	GlobVar.DB, err = sql.Open("sqlite", "server/tmp/database.db") // Updated path
	if err != nil {
		log.Fatal(err)
	}
	Migrations.Migrate()
}

// AppHandler implements http.Handler
type AppHandler struct{}

func (h AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	// Public routes
	mux.Handle("/static/", http.StripPrefix("/static", http.HandlerFunc(Handlers.HandleStatic)))
	mux.HandleFunc("/", Handlers.HandleIndex)
	mux.HandleFunc("/post/", Handlers.HandlePostPage)
	mux.HandleFunc("/Sign_In", Handlers.HandleSignIn)
	mux.HandleFunc("/Sign_Up", Handlers.HandleSignUp)
	mux.HandleFunc("/api/auth/status", Handlers.HandleAuthStatus)
	mux.HandleFunc("/api/checkEmail", Handlers.HandleIdentifierDisponibility)
	mux.HandleFunc("/api/isValidAuth", Handlers.HandleIsValidCredentials)

	// Protected routes
	mux.HandleFunc("/Comment", middleware.ValidateSession(Handlers.HandleComment))
	mux.HandleFunc("/IsLike", middleware.ValidateSession(Handlers.HandleLikeDislike))
	mux.HandleFunc("/Log_Out", middleware.ValidateSession(Handlers.HandleLogOut))
	mux.HandleFunc("/Profile_Account", middleware.ValidateSession(Handlers.HandleProfileAccount))
	mux.HandleFunc("/Update_Profile", middleware.ValidateSession(Handlers.HandleProfileUpdate))
	mux.HandleFunc("/New_Post", middleware.ValidateSession(Handlers.HandleNewPost))

	// Handle uploads
	Handlers.HandleUploads()

	mux.ServeHTTP(w, r)
}

// Handler function for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	AppHandler{}.ServeHTTP(w, r)
}

// Keep the main function for local development
func main() {
	defer GlobVar.DB.Close()

	log.Println("server start: http://localhost:8080/")
	err := http.ListenAndServe(":8080", AppHandler{})
	if err != nil {
		log.Printf("server not listener: %v", err)
	}
}