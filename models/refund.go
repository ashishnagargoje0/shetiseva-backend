package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Refund struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID   string             `bson:"order_id" json:"order_id"`
	Reason    string             `bson:"reason" json:"reason"`
	Amount    float64            `bson:"amount" json:"amount"`
	Status    string             `bson:"status" json:"status"` // e.g. Initiated, Approved
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
