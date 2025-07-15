package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client
var FriendRequests *mongo.Collection
var Users *mongo.Collection

func Init()  {
	error := godotenv.Load()
	if error != nil {
		log.Println("No .env file founded")
	}

	var uri string
	if os.Getenv("ENVIROMENT") == "development" {
		uri = os.Getenv("MONGODB_URI_DEV")
	} else {
		uri = os.Getenv("MONGODB_URI")
	}
	
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")

	MongoClient = client
	Users = MongoClient.Database("blendz").Collection("users")
	FriendRequests = MongoClient.Database("blendz").Collection("friendRequests")
}
