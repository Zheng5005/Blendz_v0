package models

import (
	"context"
	"fmt"

	"github.com/Zheng5005/Blendz_v0/db"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewFriendRequest(sender, recipient string) (*FriendRequest, error) {
  ParseUserID, err := bson.ObjectIDFromHex(sender)
	if err != nil {
		return &FriendRequest{}, fmt.Errorf("Invalid id: %w", err)
	}

  ParseRecipientID, err := bson.ObjectIDFromHex(recipient)
	if err != nil {
		return &FriendRequest{}, fmt.Errorf("Invalid id: %w", err)
	}

	friendRequest := FriendRequest{
		Sender: ParseUserID, 
		Recipient: ParseRecipientID, 
		Status: "pending", 
	}

	return &friendRequest, nil
}

func FindRequest(userId, recipientId string) (FriendRequest, error)  {
	ParseUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return FriendRequest{}, fmt.Errorf("Invalid id: %w", err)
	}

	ParseRecipientID, err := bson.ObjectIDFromHex(recipientId)
	if err != nil {
		return FriendRequest{}, fmt.Errorf("Invalid id: %w", err)
	}

	filters := []bson.M{
		{"sender": ParseUserID, "recipient": ParseRecipientID},
		{"sender": ParseRecipientID, "recipient": ParseUserID},
	}

	filter := bson.M{"$or": filters}

	var result FriendRequest
	err = db.FriendRequests.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return FriendRequest{}, fmt.Errorf("No request founded")
		}
		return FriendRequest{}, fmt.Errorf("Error: %w", err)
	}

	return result, nil
}

func InsertRequest(request FriendRequest) (bool, error)  {
	_, err := db.FriendRequests.InsertOne(context.TODO(), request)
	if err != nil {
		return false, fmt.Errorf("Failed to insert request: %w", err)
	}

	return true, nil
}

func IsFriendRequestPending(userId, recipientId string) (bool, error) {
	_, err := FindRequest(userId, recipientId)
	if err != nil {
		if err.Error() == "No request founded" {
			return false, nil
		}
		return true, fmt.Errorf("Error: %w", err)
	}

	return true, nil
}
