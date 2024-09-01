package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Following []primitive.ObjectID `bson:"following" json:"following"`
	Followers []primitive.ObjectID `bson:"followers" json:"followers"`
}
