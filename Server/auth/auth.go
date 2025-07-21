package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Zheng5005/Blendz_v0/models"
	"github.com/Zheng5005/Blendz_v0/stream"
	"github.com/Zheng5005/Blendz_v0/utils"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request)  {
	//TODO: Check the token and the newID object, cause at the moment that cookie return an error.
	// you have to use the login endpoint in order to get a valid cookie with a valid token
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

	err = stream.CreateStreamUser(newID.Hex(), user.Fullname, user.ProfilePic)
	if err != nil {
		log.Panicf("Stream user creation failed: %v", err)
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		http.Error(w, "Server Config error", http.StatusInternalServerError)
		return
	}

	token, err := utils.GenerateJWT(newID.Hex(), secret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	cookie := utils.SetCookie(token)

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Success"))
}

func Login(w http.ResponseWriter, r *http.Request)  {
	error := godotenv.Load()
	if error != nil {
		log.Println("No .env file founded")
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	} 

	var input struct {
		Email           string             `json:"email"`
		Password        string             `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	user, err := models.FindUser(input.Email)
	if err != nil {
		log.Printf("FindUser error: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		http.Error(w, "Server Config error", http.StatusInternalServerError)
		return
	}

	token, err := utils.GenerateJWT(user.ID.Hex(), secret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	cookie := utils.SetCookie(token)

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

func Logout(w http.ResponseWriter, r *http.Request)  {
	cookie := utils.ClearCookie()

	http.SetCookie(w, cookie)
	w.Write([]byte("Logout"))
}

func OnBoard(w http.ResponseWriter, r *http.Request)  {
	userId, err := utils.ParseToken(r)
	if err != nil {
		http.Error(w, "No cookie provied", http.StatusUnauthorized)
		return
	}
	//TODO: add the ProfilePic field to this handler

	var input struct {
		FullName                string             `json:"fullName"`
		BIO                     string             `json:"bio"`
		NativeLanguage          string             `json:"nativeLanguage"`
		LearningLanguage        string             `json:"learningLanguage"`
		Location                string             `json:"location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if input.FullName == "" || input.BIO == "" || input.NativeLanguage == "" || input.LearningLanguage == "" || input.Location == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	err = models.UpdateUserByID(userId, input.FullName, input.BIO, input.NativeLanguage, input.LearningLanguage, input.Location)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user, err2 := models.FindUserByID(userId)
	if err2 != nil {
		http.Error(w, "No user founded", http.StatusInternalServerError)
		return
	}

	err = stream.CreateStreamUser(user.ID.Hex(), user.Fullname, user.ProfilePic)
	if err != nil {
		log.Panicf("Stream user creation failed: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success OnBoard"))
}
