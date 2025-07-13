package stream

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	stream "github.com/GetStream/stream-chat-go/v7"
)

var Client *stream.Client

func Init()  {
	apiKey := os.Getenv("STREAM_API_KEY")
	apiSecret := os.Getenv("STREAM_API_SECRET")

	var err error
	Client, err = stream.NewClient(apiKey, apiSecret)
	if err != nil {
		log.Fatalf("Stream client init failed: %v", err)
	}
}

func CreateStreamUser(userID, fullname, profilePic string) error  {
	_, err := Client.UpsertUser(context.TODO(), &stream.User{
		ID: userID,
		Name: fullname,
		Image: profilePic,
	})

	return err
}

func GenerateStreamToken(userID string) (string, error) {
	token, err := Client.CreateToken(userID, time.Time{})
	if err != nil {
		return "", fmt.Errorf("failed to create stream token: %w", err)
	}
	return token, nil
}
