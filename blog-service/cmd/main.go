package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/omsurase/Blogging/blog-service/internal/config"
	"github.com/omsurase/Blogging/blog-service/internal/handlers"
	"github.com/omsurase/Blogging/blog-service/internal/repository"
	"github.com/omsurase/Blogging/blog-service/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
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

	// Set up gRPC connection to auth service
	log.Printf("Attempting to connect to auth service at %s", cfg.AuthServiceAddress)
	authConn, err := grpc.Dial(cfg.AuthServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	// Check the gRPC connection state
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		state := authConn.GetState()
		if state == connectivity.Ready {
			log.Println("Successfully connected to auth service")
			break
		}
		if !authConn.WaitForStateChange(ctx, state) {
			log.Fatalf("gRPC connection state did not become ready: %v", authConn.GetState())
		}
	}

	r := mux.NewRouter()

	// Initialize repository
	repo := repository.NewMongoRepository(client, cfg.MongoDB)

	// Initialize service
	blogService := service.NewBlogService(repo)

	// Initialize handlers
	blogHandler := handlers.NewBlogHandler(blogService, authConn)

	// Define routes
	r.HandleFunc("/posts", blogHandler.ValidateToken(blogHandler.CreatePost)).Methods("POST")
	r.HandleFunc("/posts", blogHandler.ValidateToken(blogHandler.GetAllPosts)).Methods("GET")
	r.HandleFunc("/posts/{id}", blogHandler.ValidateToken(blogHandler.GetPost)).Methods("GET")
	r.HandleFunc("/posts/{id}", blogHandler.ValidateToken(blogHandler.UpdatePost)).Methods("PUT")
	r.HandleFunc("/posts/{id}", blogHandler.ValidateToken(blogHandler.DeletePost)).Methods("DELETE")

	// Start server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
