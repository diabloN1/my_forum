package cookies

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate_Cookie_session generates a cryptographically secure random session ID.
func Generate_Cookie_session() (string, error) {
    id := make([]byte, 32)
    _, err := rand.Read(id)
    if err != nil {
        return "", err // Return the error instead of logging and exiting
    }
    return base64.RawStdEncoding.EncodeToString(id), nil
}

