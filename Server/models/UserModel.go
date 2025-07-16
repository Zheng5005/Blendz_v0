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

	count, err := db.Users.CountDocuments(context.TODO(), bson.D{{Key: "email", Value: user.Email}})

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

	newUser, err := db.Users.InsertOne(context.TODO(), user)
	if err != nil {
		return bson.NilObjectID, fmt.Errorf("Failed to insert user: %w", err)
	}

	return newUser.InsertedID.(bson.ObjectID), nil
}

func FindUser(email string) (UserCredentials, error) {
	filter := bson.M{"email": email}

	opts := options.FindOne().SetProjection(bson.D{
		{Key: "_id", Value: 1}, 
		{Key: "email", Value: 1}, 
		{Key: "password", Value: 1},
	})

	var result UserCredentials

	err := db.Users.FindOne(context.TODO(), filter, opts).Decode(&result)
	if err != nil {
    if err == mongo.ErrNoDocuments {
			return UserCredentials{}, fmt.Errorf("No user was found: %w", err)
		}
		return UserCredentials{}, fmt.Errorf("Error: %w", err)
	}

	return result, nil
}

func FindUserByID(id string) (User, error)  {
	ParseID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return User{}, fmt.Errorf("Invalid id: %w", err)
	}

	filter := bson.M{"_id": ParseID}

	var result User
	
	err = db.Users.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
    if err == mongo.ErrNoDocuments {
			return User{}, fmt.Errorf("No user was found: %w", err)
		}
		return User{}, fmt.Errorf("Error: %w", err)
	}

	return result, nil
}

func UpdateUserByID(id, fullname, bio, nativeLanguage, learningLanguage, location string) error  {
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

	result, err := db.Users.UpdateByID(context.TODO(), ParseID, update)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("No user was found")
	}

	return nil
}

func FindRecommendedUsers(id string) ([]User, error) {
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

	cursor, err := db.Users.Find(context.TODO(), filter)
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

	cursor, err := db.Users.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrive friends: %w", err)
	}

	var friends []User
	if err := cursor.All(context.TODO(), &friends); err != nil {
		return nil, fmt.Errorf("Failed to decode friends: %w", err)
	}

	return friends, nil
}

func AreUsersFriends(senderID, recipientID string) (bool, error) {
	// Parse both IDs
	senderObjID, err := bson.ObjectIDFromHex(senderID)
	if err != nil {
		return false, fmt.Errorf("Invalid sender ID: %w", err)
	}

	recipientObjID, err := bson.ObjectIDFromHex(recipientID)
	if err != nil {
		return false, fmt.Errorf("Invalid recipient ID: %w", err)
	}

	// Query to check if sender is in recipient's friends
	filter := bson.M{
		"_id":     recipientObjID,
		"friends": senderObjID,
	}

	// Try to find recipient with sender in friends list
	count, err := db.Users.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, fmt.Errorf("Database error: %w", err)
	}

	return count > 0, nil
}

func AddUserToFriendList(userId, newFriendId string) error {
	ParseUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("Invalid id: %w", err)
	}

	ParseFriendID, err := bson.ObjectIDFromHex(newFriendId)
	if err != nil {
		return fmt.Errorf("Invalid id: %w", err)
	}
	
	filter := bson.M{"_id": ParseUserID}
	update := bson.M{
		"$addToSet": bson.M{"friends": ParseFriendID},
	}

	result, err := db.Users.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("Failed to add friend: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("User not found")
	}

	return nil
}

func AddMutualFriends(userId, friendId string) error  {
	if err := AddUserToFriendList(userId, friendId); err != nil {
		return err
	}

	if err := AddUserToFriendList(friendId, userId); err != nil {
		return err
	}

	return nil
}

func GetFriendRequests(userID string) (map[string][]bson.M, error) {
	userObjID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// 🟡 1. Pending Requests (incoming) — you are the recipient
	pendingPipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"recipient": userObjID,
			"status":    "pending",
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "sender",
			"foreignField": "_id",
			"as":           "user",
		}}},
		{{Key: "$unwind", Value: "$user"}},
		{{Key: "$project", Value: bson.M{
			"id":     "$_id",
			"status": 1,
			"user": bson.M{
				"id":             "$user._id",
				"fullName":       "$user.fullname",
				"profilePic":     "$user.profilepic",
				"nativeLanguage": "$user.nativelanguage",
				"learningLanguage": "$user.learninglanguage",
			},
		}}},
	}

	// 🟢 2. Accepted Requests (you are sender or recipient)
	acceptedPipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"status": "accepted",
			"$or": []bson.M{
				{"sender": userObjID},
			},
		}}},
		{{Key: "$addFields", Value: bson.M{
			"otherUserId": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$sender", userObjID}},
					"$recipient",
					"$sender",
				},
			},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "otherUserId",
			"foreignField": "_id",
			"as":           "user",
		}}},
		{{Key: "$unwind", Value: "$user"}},
		{{Key: "$project", Value: bson.M{
			"id":     "$_id",
			"status": 1,
			"user": bson.M{
				"id":             "$user._id",
				"fullName":       "$user.fullname",
				"profilePic":     "$user.profilepic",
			},
		}}},
	}

	// Execute both
	pendingCursor, err := db.FriendRequests.Aggregate(context.TODO(), pendingPipeline)
	if err != nil {
		return nil, fmt.Errorf("error fetching pending: %w", err)
	}
	var pendingResults []bson.M
	if err = pendingCursor.All(context.TODO(), &pendingResults); err != nil {
		return nil, fmt.Errorf("error decoding pending: %w", err)
	}

	acceptedCursor, err := db.FriendRequests.Aggregate(context.TODO(), acceptedPipeline)
	if err != nil {
		return nil, fmt.Errorf("error fetching accepted: %w", err)
	}
	var acceptedResults []bson.M
	if err = acceptedCursor.All(context.TODO(), &acceptedResults); err != nil {
		return nil, fmt.Errorf("error decoding accepted: %w", err)
	}

	return map[string][]bson.M{
		"PendingRequest":  pendingResults,
		"AcceptedRequest": acceptedResults,
	}, nil
}

