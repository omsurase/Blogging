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

	r := mux.NewRouter()

	// Initialize repository
	repo := repository.NewMongoRepository(client, cfg.MongoDB)

	// Initialize service
	blogService := service.NewBlogService(repo)

	// Initialize handlers
	blogHandler := handlers.NewBlogHandler(blogService)

	// Define routes
	r.HandleFunc("/posts", blogHandler.CreatePost).Methods("POST")
	r.HandleFunc("/posts", blogHandler.GetAllPosts).Methods("GET")
	r.HandleFunc("/posts/{id}", blogHandler.GetPost).Methods("GET")
	r.HandleFunc("/posts/{id}", blogHandler.UpdatePost).Methods("PUT")
	r.HandleFunc("/posts/{id}", blogHandler.DeletePost).Methods("DELETE")

	// Start server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
