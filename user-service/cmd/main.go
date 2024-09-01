package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"github.com/omsurase/Blogging/user-service/internal/config"
	"github.com/omsurase/Blogging/user-service/internal/handlers"
	pb "github.com/omsurase/Blogging/user-service/internal/pb"
	"github.com/omsurase/Blogging/user-service/internal/repository"
	"github.com/omsurase/Blogging/user-service/internal/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB")
	defer client.Disconnect(ctx)

	// Initialize repository
	repo := repository.NewUserRepository(client.Database(cfg.MongoDB))

	// Initialize service
	userService := service.NewUserService(repo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Set up HTTP server
	r := mux.NewRouter()
	//r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/users/{followerID}/follow/{followeeID}", userHandler.FollowUser).Methods("POST")
	r.HandleFunc("/users/{followerID}/unfollow/{followeeID}", userHandler.UnfollowUser).Methods("POST")
	r.HandleFunc("/users/{id}/following", userHandler.GetFollowing).Methods("GET")
	r.HandleFunc("/users/{id}/followers", userHandler.GetFollowers).Methods("GET")

	// Set up gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	// Start HTTP server
	go func() {
		addr := fmt.Sprintf(":%d", cfg.ServerPort)
		log.Printf("Starting HTTP server on %s", addr)
		log.Fatal(http.ListenAndServe(addr, r))
	}()

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
	log.Printf("Starting gRPC server on :%d", cfg.GRPCPort)
	log.Fatal(grpcServer.Serve(lis))
}
