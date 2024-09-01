package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/omsurase/Blogging/auth-service/internal/config"
	"github.com/omsurase/Blogging/auth-service/internal/handlers"
	pb "github.com/omsurase/Blogging/auth-service/internal/pb"
	"github.com/omsurase/Blogging/auth-service/internal/repository"
	"github.com/omsurase/Blogging/auth-service/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type grpcServer struct {
	pb.UnimplementedAuthServiceServer
	authHandler *handlers.AuthHandler
}

func (s *grpcServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	return s.authHandler.GRPCValidateToken(ctx, req)
}

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

	userConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		state := userConn.GetState()
		if state == connectivity.Ready {
			log.Println("Successfully connected to user service")
			break
		}
		if !userConn.WaitForStateChange(ctx, state) {
			log.Fatalf("gRPC connection state did not become ready: %v", userConn.GetState())
		}
	}

	// Initialize repository
	repo := repository.NewMongoRepository(client, cfg.MongoDB)

	// Initialize service
	authService := service.NewAuthService(repo, cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, userConn)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterAuthServiceServer(s, &grpcServer{authHandler: authHandler})
		log.Printf("gRPC server listening on :%d", cfg.GRPCPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Define routes
	r.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/auth/validate", authHandler.ValidateToken).Methods("POST")

	// Start server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
