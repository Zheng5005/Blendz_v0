package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Zheng5005/Blendz_v0/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func ProtectRoute(next http.HandlerFunc) http.HandlerFunc {
	errENV := godotenv.Load()
	if errENV != nil {
		log.Println("No .env file available")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := utils.GetCokkie(r)
		if err != nil {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Invalid method")
			}
			return []byte(getEnv("JWT_SECRET", "other_key")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token or expired", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
