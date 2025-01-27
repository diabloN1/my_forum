package middleware

import (
	"context"
	Cruds "forum/Api"
	"net/http"
)

type userIdType string

func ValidateSession(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get the session cookie
        cookie, err := r.Cookie("Session_ID")
        if err != nil {
            // No session cookie found, redirect to sign-in page
            http.Redirect(w, r, "/Sign_In", http.StatusSeeOther)
            return
        }

        // Validate the session ID and get the user ID
        sessionID := cookie.Value
        userID, valid := Cruds.ValidateSessionIDAndGetUserID(sessionID)
        if !valid {
            // Invalid session, redirect to sign-in page
            http.Redirect(w, r, "/Sign_In", http.StatusSeeOther)
            return
        }
        
        // Add the user ID to the request context
        var userIdValue userIdType = "userID"
        ctx := context.WithValue(r.Context(), userIdValue, userID)
        next(w, r.WithContext(ctx))
    }
}