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

func SendFriendRequest(w http.ResponseWriter, r *http.Request)  {
	userId, err := utils.ParseToken(r)
	if err != nil {
		http.Error(w, "No cookie provied", http.StatusUnauthorized)
		return
	}

	recipientId := r.PathValue("id")
	if recipientId == "" {
		http.Error(w, "No recipientId provied", http.StatusBadRequest)
		return
	}

	// Validatind if the user isn't sending a request to themself
	if userId == recipientId {
		http.Error(w, "You can't send friend request to yourself", http.StatusBadRequest)
		return
	}

	// Validating if the user is already friends with the recipient
	isFriend, err := models.AreUsersFriends(userId, recipientId)
	if err != nil {
		log.Println("Check error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if isFriend {
		http.Error(w, "You are already friends with this user", http.StatusBadRequest)
		return
	}

	//Validating if there's a request pending
	existingRequest, err := models.IsFriendRequestPending(userId, recipientId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if existingRequest {
		http.Error(w, "A friend request already exist between you and this user", http.StatusBadRequest)
		return
	}

	newRequest, err := models.NewFriendRequest(userId, recipientId)
	if err != nil {
		log.Println("Check error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	insertedRequest, err := models.InsertRequest(*newRequest)
	if err != nil {
		log.Println("Check error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insertedRequest.Hex())
}

func AcceptFriendRequest(w http.ResponseWriter, r *http.Request)  {
	userId, err := utils.ParseToken(r)
	if err != nil {
		http.Error(w, "No cookie provied", http.StatusUnauthorized)
		return
	}

	requestId := r.PathValue("id")
	if requestId == "" {
		http.Error(w, "No requestId provied", http.StatusBadRequest)
		return
	}

	existingRequest, err := models.FindRequestByID(requestId)
	if err != nil {
		log.Println(err)
		http.Error(w, "No friend request founded", http.StatusNotFound)
		return
	}

	if existingRequest.Recipient.Hex() != userId {
		http.Error(w, "You are not authorized to accept this request", http.StatusUnauthorized)
		return
	}

	err = models.UpdateRequestByID(requestId, existingRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update the request", http.StatusInternalServerError)
		return
	}

	//Adding each user to friends list
	err = models.AddMutualFriends(userId, existingRequest.Sender.Hex())
	if err != nil {
		log.Println("Error adding friend: ", err)
		http.Error(w, "Could not add friend", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Friend request accepted"))
}

func GetFriendRequests(w http.ResponseWriter, r *http.Request)  {
	userId, err := utils.ParseToken(r)
	if err != nil {
		http.Error(w, "No cookie provied", http.StatusUnauthorized)
		return
	}

  results, err := models.GetFriendRequests(userId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
