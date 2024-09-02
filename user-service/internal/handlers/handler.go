package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/omsurase/Blogging/user-service/internal/models"
	pb "github.com/omsurase/Blogging/user-service/internal/pb"
	"github.com/omsurase/Blogging/user-service/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("CreateUser: Received request to create user with ID: %s", req.UserId)
	fmt.Println(req.UserId)
	objectID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		log.Printf("CreateUser: Error converting user ID to ObjectID. User ID: %s, Error: %v", req.UserId, err)
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	user := &models.User{
		ID: objectID,
	}

	log.Printf("CreateUser: Attempting to create user with ObjectID: %s", objectID.Hex())
	err = h.userService.CreateUser(user)
	if err != nil {
		log.Printf("CreateUser: Failed to create user. User ID: %s, Error: %v", objectID.Hex(), err)
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	log.Printf("CreateUser: Successfully created user with ID: %s", objectID.Hex())
	return &pb.CreateUserResponse{
		Success: true,
		Message: fmt.Sprintf("User created successfully with ID: %s", objectID.Hex()),
	}, nil
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	log.Printf("GetUser: Received request to get user with ID: %s", userID)

	user, err := h.userService.GetUser(userID)
	if err != nil {
		log.Printf("GetUser: Error retrieving user with ID %s: %v", userID, err)
		http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusNotFound)
		return
	}

	log.Printf("GetUser: Successfully retrieved user with ID: %s", userID)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followerID := vars["followerID"]
	followeeID := vars["followeeID"]
	log.Printf("FollowUser: Received request for user %s to follow user %s", followerID, followeeID)

	err := h.userService.FollowUser(followerID, followeeID)
	if err != nil {
		log.Printf("FollowUser: Error following user. Follower ID: %s, Followee ID: %s, Error: %v", followerID, followeeID, err)
		http.Error(w, fmt.Sprintf("Failed to follow user: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("FollowUser: User %s successfully followed user %s", followerID, followeeID)
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followerID := vars["followerID"]
	followeeID := vars["followeeID"]
	log.Printf("UnfollowUser: Received request for user %s to unfollow user %s", followerID, followeeID)

	err := h.userService.UnfollowUser(followerID, followeeID)
	if err != nil {
		log.Printf("UnfollowUser: Error unfollowing user. Follower ID: %s, Followee ID: %s, Error: %v", followerID, followeeID, err)
		http.Error(w, fmt.Sprintf("Failed to unfollow user: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("UnfollowUser: User %s successfully unfollowed user %s", followerID, followeeID)
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	log.Printf("GetFollowing: Received request to get users followed by user with ID: %s", userID)

	following, err := h.userService.GetFollowing(userID)
	if err != nil {
		log.Printf("GetFollowing: Error retrieving following list for user %s: %v", userID, err)
		http.Error(w, fmt.Sprintf("Failed to get following list: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("GetFollowing: Successfully retrieved following list for user %s. Count: %d", userID, len(following))
	json.NewEncoder(w).Encode(following)
}

func (h *UserHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	log.Printf("GetFollowers: Received request to get followers of user with ID: %s", userID)

	followers, err := h.userService.GetFollowers(userID)
	if err != nil {
		log.Printf("GetFollowers: Error retrieving followers list for user %s: %v", userID, err)
		http.Error(w, fmt.Sprintf("Failed to get followers list: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("GetFollowers: Successfully retrieved followers list for user %s. Count: %d", userID, len(followers))
	json.NewEncoder(w).Encode(followers)
}
