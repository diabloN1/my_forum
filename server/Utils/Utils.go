package Utils

import (
	Cruds "forum/Api"
	"net/http"
)

func GetCurrentUserId(r *http.Request) string{
	// Retrieve the user ID from the cookie
	var isValid bool
	cookie, err := r.Cookie("Session_ID")
	if err != nil {
		return ""
	}
	userID := ""
	if cookie != nil {
		// Validate the session ID and get the user ID
		sessionID := cookie.Value
		userID, isValid = Cruds.ValidateSessionIDAndGetUserID(sessionID)
		if !isValid {
			return ""
		}
	}
	return userID
}
