package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Zheng5005/Blendz_v0/models"
	"github.com/Zheng5005/Blendz_v0/utils"
	"github.com/joho/godotenv"
)

func Signup(w http.ResponseWriter, r *http.Request)  {
	error := godotenv.Load()
	if error != nil {
		log.Println("No .env file founded")
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	} 

	var input struct {
		Fullname        string             `json:"fullName"`
		Email           string             `json:"email"`
		Password        string             `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	user := models.NewUser(input.Fullname, input.Email, input.Password)

	newID, err := models.InsertUser(*user)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	str := fmt.Sprintf("%v", newID)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		http.Error(w, "Server Config error", http.StatusInternalServerError)
		return
	}

	token, err := utils.GenerateJWT(str, secret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	cookie := utils.SetCookie(token)

	// TODO: CREATE USER IN STREAM 

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Success"))
}

func Login(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("Login"))
}

func Logout(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("Logout"))
}
