package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SupportTicket struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId" binding:"required"`
	Subject   string             `bson:"subject" json:"subject" binding:"required"`
	Message   string             `bson:"message" json:"message" binding:"required"`
	Status    string             `bson:"status" json:"status,omitempty"`
	CreatedAt int64              `bson:"created_at" json:"createdAt,omitempty"`
	UpdatedAt int64              `bson:"updated_at" json:"updatedAt,omitempty"`
}
