package repository

import (
	"context"
	"fmt"

	"github.com/omsurase/Blogging/blog-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	client *mongo.Client
	dbName string
}

func NewMongoRepository(client *mongo.Client, dbName string) *MongoRepository {
	return &MongoRepository{
		client: client,
		dbName: dbName,
	}
}

func (r *MongoRepository) CreatePost(post *models.Post) error {
	collection := r.client.Database(r.dbName).Collection("posts")
	_, err := collection.InsertOne(context.Background(), post)
	return err
}

func (r *MongoRepository) GetAllPosts() ([]models.Post, error) {
	collection := r.client.Database(r.dbName).Collection("posts")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var posts []models.Post
	if err = cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *MongoRepository) GetPost(id string) (*models.Post, error) {
	collection := r.client.Database(r.dbName).Collection("posts")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %v", err)
	}

	var post models.Post
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("post not found: %s", id)
		}
		return nil, fmt.Errorf("error fetching post: %v", err)
	}
	return &post, nil
}

func (r *MongoRepository) UpdatePost(id string, post *models.Post) error {
	collection := r.client.Database(r.dbName).Collection("posts")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": post},
	)
	return err
}

func (r *MongoRepository) DeletePost(id string) error {
	collection := r.client.Database(r.dbName).Collection("posts")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}
