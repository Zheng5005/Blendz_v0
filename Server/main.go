package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Zheng5005/Blendz_v0/auth"
	"github.com/Zheng5005/Blendz_v0/db"
	"github.com/Zheng5005/Blendz_v0/stream"
	"github.com/joho/godotenv"
)

func main()  {
	error := godotenv.Load()
	if error != nil {
		log.Println("No .env file founded")
	}

	// MongoDB connection
	db.Init()

	// getStream connection
	stream.Init()

	//Disconnect MongoDB when main() exits
	defer func ()  {
		if err := db.MongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatal("Failed to Disconnect MongoDB: ", err)
		}
	}()

	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("POST /api/auth/signup", auth.Signup)
	mux.HandleFunc("POST /api/auth/login", auth.Login)
	mux.HandleFunc("POST /api/auth/logout", auth.Logout)
	// http.HandleFunc("GET /posts/{id}", handlePost2)

	s, err := makeServer(mux)
	if err != nil {
		log.Fatalf("Couldn't make a new server")
	}

	log.Println("Server running...")
	log.Fatal(s.ListenAndServe())
}

func makeServer(mux *http.ServeMux) (s *http.Server, err error)  {
	error := godotenv.Load()
	if error != nil {
		log.Println("No .env file founded")
		return nil, error
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = ":8081"
	}

	server := &http.Server{
		Addr: port,
		Handler: mux,
	}

	return server, nil
}
