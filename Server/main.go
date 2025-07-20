package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Zheng5005/Blendz_v0/auth"
	"github.com/Zheng5005/Blendz_v0/chat"
	"github.com/Zheng5005/Blendz_v0/db"
	"github.com/Zheng5005/Blendz_v0/middleware"
	"github.com/Zheng5005/Blendz_v0/stream"
	"github.com/Zheng5005/Blendz_v0/user"
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
	mux.HandleFunc("POST /api/auth/onboarding", middleware.ProtectRoute(auth.OnBoard))

	// Users routes
	mux.HandleFunc("GET /api/users", middleware.ProtectRoute(user.GetRecommendedUsers))
	mux.HandleFunc("GET /api/users/me", middleware.ProtectRoute(user.GetMeAuth))
	mux.HandleFunc("GET /api/users/friends", middleware.ProtectRoute(user.GetMyFriends))
	mux.HandleFunc("POST /api/users/friend-request/{id}", middleware.ProtectRoute(user.SendFriendRequest))
	mux.HandleFunc("PUT /api/users/friend-request/{id}/accept", middleware.ProtectRoute(user.AcceptFriendRequest))
	mux.HandleFunc("GET /api/users/friend-requests", middleware.ProtectRoute(user.GetFriendRequests))
	mux.HandleFunc("GET /api/users/outgoing-friend-requests", middleware.ProtectRoute(user.GetOutgoingFriendRequest))

	//Chat routes
	mux.HandleFunc("GET /api/chat/token", middleware.ProtectRoute(chat.GetStreamToken))

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

	muxWithCors := middleware.CorsMiddleware(mux)

	server := &http.Server{
		Addr: port,
		Handler: muxWithCors,
	}

	return server, nil
}
