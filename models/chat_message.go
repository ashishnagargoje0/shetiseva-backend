package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatMessage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Message   string             `bson:"message" json:"message"`
	Reply     string             `bson:"reply" json:"reply"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

type ChatMessageInput struct {
	Message string `json:"message" binding:"required"`
}
