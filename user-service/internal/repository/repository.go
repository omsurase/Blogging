package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/omsurase/Blogging/user-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	log.Println("Initializing UserRepository")
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Attempting to create user: %+v", user)
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}
	log.Printf("User created successfully. Inserted ID: %v", result.InsertedID)
	return nil
}

func (r *UserRepository) GetByID(id primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Attempting to fetch user with ID: %s", id.Hex())
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("User not found with ID: %s", id.Hex())
			return nil, fmt.Errorf("user not found: %w", err)
		}
		log.Printf("Error fetching user with ID %s: %v", id.Hex(), err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	log.Printf("Successfully fetched user with ID: %s", id.Hex())
	return &user, nil
}

func (r *UserRepository) Follow(followerID, followeeID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Attempting to make user %s follow user %s", followerID.Hex(), followeeID.Hex())

	// Update follower's following list
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": followerID},
		bson.M{"$addToSet": bson.M{"following": followeeID}},
	)
	if err != nil {
		log.Printf("Error updating follower %s: %v", followerID.Hex(), err)
		return fmt.Errorf("failed to update follower: %w", err)
	}

	// Update followee's followers list
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": followeeID},
		bson.M{"$addToSet": bson.M{"followers": followerID}},
	)
	if err != nil {
		log.Printf("Error updating followee %s: %v", followeeID.Hex(), err)
		return fmt.Errorf("failed to update followee: %w", err)
	}

	log.Printf("Successfully made user %s follow user %s", followerID.Hex(), followeeID.Hex())
	return nil
}

func (r *UserRepository) Unfollow(followerID, followeeID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Attempting to make user %s unfollow user %s", followerID.Hex(), followeeID.Hex())

	// Update follower's following list
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": followerID},
		bson.M{"$pull": bson.M{"following": followeeID}},
	)
	if err != nil {
		log.Printf("Error updating follower %s: %v", followerID.Hex(), err)
		return fmt.Errorf("failed to update follower: %w", err)
	}

	// Update followee's followers list
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": followeeID},
		bson.M{"$pull": bson.M{"followers": followerID}},
	)
	if err != nil {
		log.Printf("Error updating followee %s: %v", followeeID.Hex(), err)
		return fmt.Errorf("failed to update followee: %w", err)
	}

	log.Printf("Successfully made user %s unfollow user %s", followerID.Hex(), followeeID.Hex())
	return nil
}

func (r *UserRepository) GetFollowing(userID primitive.ObjectID) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Attempting to fetch users followed by user %s", userID.Hex())

	user, err := r.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	var following []*models.User
	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": user.Following}})
	if err != nil {
		log.Printf("Error fetching following users for user %s: %v", userID.Hex(), err)
		return nil, fmt.Errorf("failed to fetch following users: %w", err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &following)
	if err != nil {
		log.Printf("Error decoding following users for user %s: %v", userID.Hex(), err)
		return nil, fmt.Errorf("failed to decode following users: %w", err)
	}

	log.Printf("Successfully fetched %d users followed by user %s", len(following), userID.Hex())
	return following, nil
}

func (r *UserRepository) GetFollowers(userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Attempting to fetch followers of user %s", userID.Hex())

	var user struct {
		Followers []primitive.ObjectID `bson:"followers"`
	}
	err := r.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		log.Printf("Error fetching user %s: %v", userID.Hex(), err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	log.Printf("Successfully fetched %d followers of user %s", len(user.Followers), userID.Hex())
	return user.Followers, nil
}
