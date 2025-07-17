package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID    string             `bson:"order_id" json:"order_id"`
	UserID     string             `bson:"user_id" json:"user_id"`
	Method     string             `bson:"method" json:"method"` // e.g., UPI, Card
	Amount     float64            `bson:"amount" json:"amount"`
	Status     string             `bson:"status" json:"status"` // Paid, Failed, Refunded
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
