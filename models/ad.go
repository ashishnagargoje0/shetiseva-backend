package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MarketplaceAd struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"`
	Price       float64            `bson:"price" json:"price"`
	ImageURL    string             `bson:"image_url" json:"image_url"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}
