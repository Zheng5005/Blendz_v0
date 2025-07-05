package models

import (
	"context"
	"fmt"

	"github.com/Zheng5005/Blendz_v0/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Fullname        string             `json:"fullName"`
	Email           string             `json:"email"`
	Password        string             `json:"password"`
	BIO             string             `json:"bio"`
	ProfilePic      string             `json:"profilePic"`
	NativeLanguage  string             `json:"nativeLanguage"`
	LearningLanguage string            `json:"learningLanguage"`
	Location        string             `json:"location"`
	IsOnboarded     bool               `json:"isOnboarded"`
	Friends         []int              `json:"friends"`
}

func ValidateUser(user User) error  {
	if user.Fullname == "" || user.Email == "" || user.Password == "" {
		return fmt.Errorf("Must have all required fields")
	}

	if len(user.Password) < 6 {
		return fmt.Errorf("Password has to be at least 6 characters")
	}

	return nil
}

func GenerateHashedPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error hashing password: %w", err)
	}

	return string(bytes), nil
}

func InsertUser(user User) error {
	if err := ValidateUser(user); err != nil {
		return err
	}

	collection := db.MongoClient.Database(db.DB).Collection("users")

	// Before saving the user, it needs to hashed the password
	hashedPassword, err := GenerateHashedPassword(user.Password)
	if err != nil {
		return fmt.Errorf("Failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	user.ID = primitive.NewObjectID()

	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return fmt.Errorf("Failed to hash password: %w", err)
	}

	return nil
}
