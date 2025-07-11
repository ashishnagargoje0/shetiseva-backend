package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"` // Add this for order ID
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items     []CartItem         `bson:"items" json:"items"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
