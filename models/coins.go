package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// CoinsBalance represents the coin balance of a user
type CoinsBalance struct {
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Coins  int                `bson:"coins" json:"coins"`
}
