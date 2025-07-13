package models

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/Zheng5005/Blendz_v0/db"
	"github.com/Zheng5005/Blendz_v0/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID              bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Fullname        string             `json:"fullName" bson:"fullname"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"password" bson:"password"`
	BIO             string             `json:"bio" bson:"bio"`
	ProfilePic      string             `json:"profilePic" bson:"profilepic"`
	NativeLanguage  string             `json:"nativeLanguage" bson:"nativelanguage"`
	LearningLanguage string            `json:"learningLanguage" bson:"learninglanguage"`
	Location        string             `json:"location" bson:"location"`
	IsOnboarded     bool               `json:"isOnboarded" bson:"isonboarded"`
	Friends         []bson.ObjectID `json:"friends,omitempty" bson:"friends,omitempty"`
}

type UserCredentials struct {
	ID              bson.ObjectID `bson:"_id,omitempty"`
	Email           string             `bson:"email"`
	Password        string             `bson:"password"`
}

func NewUser(fullName string, email string, password string) *User {
	user := User{
		Fullname: fullName, 
		Email: email, 
		Password: password,
		ProfilePic: "https://avatar.iran.liara.run/public/12.png",
		Friends: []bson.ObjectID{},
	}

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

func InsertUser(user User) (bson.ObjectID, error) {
	if err := ValidateUser(user); err != nil {
		log.Print(err)
		return bson.NilObjectID, err
	}

	// Before saving the user, it needs to hashed the password
	hashedPassword, err := utils.GenerateHashedPassword(user.Password)
	if err != nil {
		log.Print(err)
		return bson.NilObjectID, fmt.Errorf("Failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	if user.ProfilePic == "" {
		user.ProfilePic = "https://avatar.iran.liara.run/public/12.png"
	}
	user.ID = bson.NewObjectID()

	var collection = db.MongoClient.Database(db.DB).Collection("users")

	newUser, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return bson.NilObjectID, fmt.Errorf("Failed to insert user: %w", err)
	}

	return newUser.InsertedID.(bson.ObjectID), nil
}

func FindUser(email string) (UserCredentials, error) {
	collection := db.MongoClient.Database(db.DB).Collection("users")
	filter := bson.M{"email": email}

	opts := options.FindOne().SetProjection(bson.D{
		{Key: "_id", Value: 1}, 
		{Key: "email", Value: 1}, 
		{Key: "password", Value: 1},
	})

	var result UserCredentials

	err := collection.FindOne(context.TODO(), filter, opts).Decode(&result)
	if err != nil {
    if err == mongo.ErrNoDocuments {
			return UserCredentials{}, fmt.Errorf("No user was found: %w", err)
		}
		return UserCredentials{}, fmt.Errorf("Error: %w", err)
	}

	return result, nil
}
