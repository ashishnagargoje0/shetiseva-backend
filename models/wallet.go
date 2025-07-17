package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletTransaction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Amount    float64            `bson:"amount" json:"amount"`
	Type      string             `bson:"type" json:"type"` // Credit or Debit
	Reason    string             `bson:"reason" json:"reason"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
