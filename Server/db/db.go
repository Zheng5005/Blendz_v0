package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Init()  {
	error := godotenv.Load()
	if error != nil {
		log.Println("No .env file founded")
	}

	uri := os.Getenv("MONGODB_URI")
	
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func ()  {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}
