package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/omsurase/Blogging/auth-service/internal/config"
	"github.com/omsurase/Blogging/auth-service/internal/handlers"
	"github.com/omsurase/Blogging/auth-service/internal/repository"
	"github.com/omsurase/Blogging/auth-service/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	//log.Printf("%s", cfg)
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

	r := mux.NewRouter()

	// Initialize repository
	repo := repository.NewMongoRepository(client, cfg.MongoDB)

	// Initialize service
	authService := service.NewAuthService(repo, cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Define routes
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/validate", authHandler.ValidateToken).Methods("POST")

	// Start server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
