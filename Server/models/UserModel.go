package models

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/Zheng5005/Blendz_v0/db"
	"github.com/Zheng5005/Blendz_v0/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func NewUser(fullName string, email string, password string) *User {
	user := User{Fullname: fullName, Email: email, Password: password}

	return &user
}

func ValidateUser(user User) error  {
	if user.Fullname == "" || user.Email == "" || user.Password == "" {
		return fmt.Errorf("Must have all required fields")
	}

	if len(user.Password) < 6 {
		return fmt.Errorf("Password has to be at least 6 characters")
	}

	emailRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`

	// Compile the regex
	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(user.Email) {
		return fmt.Errorf("Invalid email")
	}

	var collection = db.MongoClient.Database(db.DB).Collection("users")

	count, err := collection.CountDocuments(context.TODO(), bson.D{{Key: "email", Value: user.Email}})

	if err != nil {
		return fmt.Errorf("Database error: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("Email already registered")
	}

	return nil
}

func InsertUser(user User) (*mongo.InsertOneResult, error) {
	if err := ValidateUser(user); err != nil {
		log.Print(err)
		return nil, err
	}

	// Before saving the user, it needs to hashed the password
	hashedPassword, err := utils.GenerateHashedPassword(user.Password)
	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("Failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	if user.ProfilePic == "" {
		user.ProfilePic = "https://avatar.iran.liara.run/public/12.png"
	}

	var collection = db.MongoClient.Database(db.DB).Collection("users")

	newUser, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, fmt.Errorf("Failed to insert user: %w", err)
	}

	return newUser, nil
}
