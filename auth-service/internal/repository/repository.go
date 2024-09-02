package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/omsurase/Blogging/auth-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByUsername(username string) (*models.User, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(client *mongo.Client, dbName string) *MongoRepository {
	collection := client.Database(dbName).Collection("users")
	return &MongoRepository{collection: collection}
}

func (r *MongoRepository) Create(user *models.User) error {
	objectID := primitive.NewObjectID()

	doc := bson.M{
		"_id":      objectID,
		"username": user.Username,
		"password": user.Password,
	}

	_, err := r.collection.InsertOne(context.Background(), doc)
	if err != nil {
		return err
	}

	user.ID = objectID.Hex()

	return nil
}

func (r *MongoRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	fmt.Println(user.ID)

	return &user, nil
}
