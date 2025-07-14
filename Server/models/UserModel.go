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

func NewUser(fullName string, email string, password string) *User {
	//TODO: Making a random number for the photo
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
		//TODO: Making a random number for the photo
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

func FindUserByID(id string) (User, error)  {
	collection := db.MongoClient.Database(db.DB).Collection("users")
	ParseID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return User{}, fmt.Errorf("Invalid id: %w", err)
	}

	filter := bson.M{"_id": ParseID}

	var result User
	
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
    if err == mongo.ErrNoDocuments {
			return User{}, fmt.Errorf("No user was found: %w", err)
		}
		return User{}, fmt.Errorf("Error: %w", err)
	}

	return result, nil
}

func UpdateUserByID(id, fullname, bio, nativeLanguage, learningLanguage, location string) error  {
	collection := db.MongoClient.Database(db.DB).Collection("users")
	ParseID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Invalid id: %w", err)
	}

	updatingUser := OnBoardingUser{
		Fullname: fullname, 
		BIO: bio, 
		NativeLanguage: nativeLanguage, 
		LearningLanguage: learningLanguage, 
		Location: location, 
		IsOnboarded: true,
	}

	update := bson.M{
		"$set": updatingUser,
	}

	result, err := collection.UpdateByID(context.TODO(), ParseID, update)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("No user was found")
	}

	return nil
}

func FindRecommendedUsers(id string) ([]User, error) {
	collection := db.MongoClient.Database(db.DB).Collection("users")
	ParseID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("Invalid id: %w", err)
	}

	user, err := FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("No user Found: %w", err)
	}

  filters := []bson.M{
		{"_id": bson.M{"$ne": ParseID}},
		{"isonboarded": true},
	}

	if len(user.Friends) > 0 {
		filters = append(filters, bson.M{
			"_id": bson.M{"$nin": user.Friends},
		})
	}

	filter := bson.M{"$and": filters}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("Error retriving documents: %w", err)
	}

	var result []User
	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, fmt.Errorf("Error decoding documents: %w", err)
	}

	return result, nil
}

func GetFriends(id string) ([]User, error)  {
	collection := db.MongoClient.Database(db.DB).Collection("users")
	user, err := FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("No user Found: %w", err)
	}

	if len(user.Friends) == 0 {
		return []User{}, nil
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": user.Friends,
		},
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrive friends: %w", err)
	}

	var friends []User
	if err := cursor.All(context.TODO(), &friends); err != nil {
		return nil, fmt.Errorf("Failed to decode friends: %w", err)
	}

	return friends, nil
}
