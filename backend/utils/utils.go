package utils

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Hash a plain password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// Compare a plain password with a hashed password
func CheckPassword(hashedPassword, plainPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
    return err == nil
}

// Generate a random session ID
func NewSessionID() string {
    bytes := make([]byte, 32)
    _, err := rand.Read(bytes)
    if err != nil {
        return ""
    }
    return hex.EncodeToString(bytes)
}

// Name of the cookie used for sessions
const SessionCookieName = "session_id"

// SetSessionCookie writes a secure session cookie to the response
func SetSessionCookie(w http.ResponseWriter, sessionID string) {
    http.SetCookie(w, &http.Cookie{
        Name:     SessionCookieName,
        Value:    sessionID,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // set to true if using HTTPS
        SameSite: http.SameSiteLaxMode,
        Expires:  time.Now().Add(7 * 24 * time.Hour),
    })
}

// GetSessionIDFromRequest extracts the session cookie from the request
func GetSessionIDFromRequest(r *http.Request) (string, error) {
    cookie, err := r.Cookie(SessionCookieName)
    if err != nil {
        return "", err
    }
    return cookie.Value, nil
}

// ClearSessionCookie removes the session cookie
func ClearSessionCookie(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:     SessionCookieName,
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteLaxMode,
        Expires:  time.Unix(0, 0),
    })
}

/*package utils

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// Compare a plain password with a hashed password
func CheckPassword(hashedPassword, plainPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
    return err == nil
}

// Generate a random session ID
func NewSessionID() string {
    bytes := make([]byte, 32)
    _, err := rand.Read(bytes)
    if err != nil {
        return ""
    }
    return hex.EncodeToString(bytes)
}*/

/*func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}*/
