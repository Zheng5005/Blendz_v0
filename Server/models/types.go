package models

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID              bson.ObjectID      `json:"_id,omitempty" bson:"_id,omitempty"`
	Fullname        string             `json:"fullName" bson:"fullname"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"password" bson:"password"`
	BIO             string             `json:"bio" bson:"bio"`
	ProfilePic      string             `json:"profilePic" bson:"profilepic"`
	NativeLanguage  string             `json:"nativeLanguage" bson:"nativelanguage"`
	LearningLanguage string            `json:"learningLanguage" bson:"learninglanguage"`
	Location        string             `json:"location" bson:"location"`
	IsOnboarded     bool               `json:"isOnboarded" bson:"isonboarded"`
	Friends         []bson.ObjectID    `json:"friends,omitempty" bson:"friends,omitempty"`
}

type UserCredentials struct {
	ID              bson.ObjectID      `bson:"_id,omitempty"`
	Email           string             `bson:"email"`
	Password        string             `bson:"password"`
}

type OnBoardingUser struct {
	Fullname        string        `json:"fullName" bson:"fullname"`
	BIO             string        `json:"bio" bson:"bio"`
	NativeLanguage  string        `json:"nativeLanguage" bson:"nativelanguage"`
	LearningLanguage string       `json:"learningLanguage" bson:"learninglanguage"`
	Location        string        `json:"location" bson:"location"`
	IsOnboarded     bool          `json:"isOnboarded" bson:"isonboarded"`
}

type FriendRequest struct {
	Sender       bson.ObjectID       `json:"sender" bson:"sender"`
	Recipient    bson.ObjectID       `json:"recipient" bson:"recipient"`
	Status       string              `json:"status" bson:"status"`
}
