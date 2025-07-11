package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReturnRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrderID   primitive.ObjectID `bson:"order_id" json:"orderId" binding:"required"`
	Reason    string             `bson:"reason" json:"reason" binding:"required"`
	Status    string             `bson:"status" json:"status,omitempty"`  // ✅ Add this line
	CreatedAt int64              `bson:"created_at" json:"createdAt,omitempty"`
}

type RefundRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrderID   primitive.ObjectID `bson:"order_id" json:"orderId" binding:"required"`
	Amount    float64            `bson:"amount" json:"amount" binding:"required"`
	Reason    string             `bson:"reason" json:"reason" binding:"required"`
	Status    string             `bson:"status" json:"status,omitempty"`
	CreatedAt int64              `bson:"created_at" json:"createdAt,omitempty"`
}
