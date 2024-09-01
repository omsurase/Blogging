package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/omsurase/Blogging/user-service/internal/models"
	"github.com/omsurase/Blogging/user-service/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateUser: Received request to create a new user")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("CreateUser: Error decoding request body: %v", err)
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	err = h.userService.CreateUser(&user)
	if err != nil {
		log.Printf("CreateUser: Error creating user: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("CreateUser: Successfully created user with ID: %s", user.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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
