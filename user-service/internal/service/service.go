package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/omsurase/Blogging/user-service/internal/models"
	"github.com/omsurase/Blogging/user-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	log.Println("Initializing UserService")
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *models.User) error {
	log.Printf("CreateUser: Attempting to create user with ID: %s", user.ID)
	err := s.repo.Create(user)
	if err != nil {
		log.Printf("CreateUser: Failed to create user. Error: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}
	log.Printf("CreateUser: Successfully created user with ID: %s", user.ID.Hex())
	return nil
}

func (s *UserService) GetUser(id string) (*models.User, error) {
	log.Printf("GetUser: Attempting to fetch user with ID: %s", id)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("GetUser: Invalid user ID format. Error: %v", err)
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}
	user, err := s.repo.GetByID(objectID)
	if err != nil {
		log.Printf("GetUser: Failed to fetch user. Error: %v", err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	log.Printf("GetUser: Successfully fetched user with ID: %s", id)
	return user, nil
}

func (s *UserService) FollowUser(followerID, followeeID string) error {
	log.Printf("FollowUser: User %s attempting to follow user %s", followerID, followeeID)
	followerObjectID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		log.Printf("FollowUser: Invalid follower ID format. Error: %v", err)
		return fmt.Errorf("invalid follower ID format: %w", err)
	}
	followeeObjectID, err := primitive.ObjectIDFromHex(followeeID)
	if err != nil {
		log.Printf("FollowUser: Invalid followee ID format. Error: %v", err)
		return fmt.Errorf("invalid followee ID format: %w", err)
	}
	if followerObjectID == followeeObjectID {
		log.Printf("FollowUser: User %s attempted to follow themselves", followerID)
		return errors.New("user cannot follow themselves")
	}
	err = s.repo.Follow(followerObjectID, followeeObjectID)
	if err != nil {
		log.Printf("FollowUser: Failed to follow user. Error: %v", err)
		return fmt.Errorf("failed to follow user: %w", err)
	}
	log.Printf("FollowUser: User %s successfully followed user %s", followerID, followeeID)
	return nil
}

func (s *UserService) UnfollowUser(followerID, followeeID string) error {
	log.Printf("UnfollowUser: User %s attempting to unfollow user %s", followerID, followeeID)
	followerObjectID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		log.Printf("UnfollowUser: Invalid follower ID format. Error: %v", err)
		return fmt.Errorf("invalid follower ID format: %w", err)
	}
	followeeObjectID, err := primitive.ObjectIDFromHex(followeeID)
	if err != nil {
		log.Printf("UnfollowUser: Invalid followee ID format. Error: %v", err)
		return fmt.Errorf("invalid followee ID format: %w", err)
	}
	err = s.repo.Unfollow(followerObjectID, followeeObjectID)
	if err != nil {
		log.Printf("UnfollowUser: Failed to unfollow user. Error: %v", err)
		return fmt.Errorf("failed to unfollow user: %w", err)
	}
	log.Printf("UnfollowUser: User %s successfully unfollowed user %s", followerID, followeeID)
	return nil
}

func (s *UserService) GetFollowing(userID string) ([]*models.User, error) {
	log.Printf("GetFollowing: Fetching users followed by user %s", userID)
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("GetFollowing: Invalid user ID format. Error: %v", err)
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}
	following, err := s.repo.GetFollowing(objectID)
	if err != nil {
		log.Printf("GetFollowing: Failed to fetch following users. Error: %v", err)
		return nil, fmt.Errorf("failed to fetch following users: %w", err)
	}
	log.Printf("GetFollowing: Successfully fetched %d users followed by user %s", len(following), userID)
	return following, nil
}

func (s *UserService) GetFollowers(userID string) ([]*models.User, error) {
	log.Printf("GetFollowers: Fetching followers of user %s", userID)
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("GetFollowers: Invalid user ID format. Error: %v", err)
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}
	followers, err := s.repo.GetFollowers(objectID)
	if err != nil {
		log.Printf("GetFollowers: Failed to fetch followers. Error: %v", err)
		return nil, fmt.Errorf("failed to fetch followers: %w", err)
	}
	log.Printf("GetFollowers: Successfully fetched %d followers of user %s", len(followers), userID)
	return followers, nil
}
