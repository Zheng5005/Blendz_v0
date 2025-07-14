package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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

func ParseToken(r *http.Request) (string, error)  {
	errENV := godotenv.Load()
	if errENV != nil {
		log.Println("No .env file available")
	}

	cookie, err := r.Cookie("Blendz_Session")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("No cookie found: %w", err)
		}
		return "", fmt.Errorf("Error reading cookie: %w", err)
	}

	tokenStr := cookie.Value
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid method")
		}
		return []byte(getEnv("JWT_SECRET", "other_key")), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("Invalid or expired token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invalid claims")
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		return "", fmt.Errorf("userId missing")
	}

	return userID, nil
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

func ClearCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:    "Blendz_Session",
		Value:   "",
		MaxAge: -1,
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

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
