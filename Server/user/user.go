package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Zheng5005/Blendz_v0/models"
	"github.com/Zheng5005/Blendz_v0/utils"
)

func GetRecommendedUsers(w http.ResponseWriter, r *http.Request)  {
	userId, err := utils.ParseToken(r)
	if err != nil {
		http.Error(w, "No cookie provied", http.StatusUnauthorized)
		return
	}

	recommendedUsers, err := models.FindRecommendedUsers(userId)
	if err != nil {
		log.Println(err)
		http.Error(w, "No users found", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendedUsers)
}

func GetMyFriends(w http.ResponseWriter, r *http.Request)  {
	userId, err := utils.ParseToken(r)
	if err != nil {
		http.Error(w, "No cookie provied", http.StatusUnauthorized)
		return
	}

	friends, err := models.GetFriends(userId)
	if err != nil {
		log.Println(err)
		http.Error(w, "No users found", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}
