package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashedPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error hashing password: %w", err)
	}

	return string(bytes), nil
}

func GenerateJWT(userID, secret string) (string, error)  {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func SetCookie(token string) *http.Cookie {
	cookie := &http.Cookie{
		Name:    "Blendz_Session",
		Value:   token,
		Expires: time.Now().Add(168 * time.Hour), // Cookie expires in 7 days
		Path:    "/",                            // Valid for all paths
		HttpOnly: true,                           // Not accessible via client-side scripts
		Secure:  true,                           // Only sent over HTTPS
		SameSite: http.SameSiteLaxMode,               // Controls cross-site requests
	}

	return cookie
}

func GetCokkie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("Blendz_Session")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, fmt.Errorf("No cookie found: %w", err)
		}
		return nil, fmt.Errorf("Error reading cookie: %w", err)
	}

	return cookie, nil
}
