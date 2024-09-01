package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/omsurase/Blogging/auth-service/internal/models"
	pb "github.com/omsurase/Blogging/auth-service/internal/pb"
	userpb "github.com/omsurase/Blogging/auth-service/internal/pb/userpb"
	"github.com/omsurase/Blogging/auth-service/internal/service"
	"google.golang.org/grpc"
)

type AuthHandler struct {
	authService *service.AuthService
	userClient  userpb.UserServiceClient
}

func NewAuthHandler(authService *service.AuthService, userConn *grpc.ClientConn) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userClient:  userpb.NewUserServiceClient(userConn),
	}
}

// GRPCValidateToken handles token validation for gRPC requests
func (h *AuthHandler) GRPCValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	valid, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	return &pb.ValidateTokenResponse{
		Valid: valid,
	}, nil
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	log.Printf("Register handler called")
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("request received.")
	fmt.Println((user))
	userId, token, err := h.authService.Register(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createUserReq := &userpb.CreateUserRequest{
		UserId: userId,
	}
	_, err = h.userClient.CreateUser(context.Background(), createUserReq)
	if err != nil {
		log.Printf("Failed to create user in user microservice: %v", err)
		http.Error(w, "Failed to complete user registration", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully", "token": token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("request recieved.")

	token, err := h.authService.Login(credentials.Username, credentials.Password)

	if err != nil {
		log.Printf("request recieved.2")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var tokenReq models.TokenValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&tokenReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	valid, err := h.authService.ValidateToken(tokenReq.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"valid": valid})
}
