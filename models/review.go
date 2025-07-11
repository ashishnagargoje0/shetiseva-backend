package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductReview struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`
	ProductID primitive.ObjectID `bson:"product_id" json:"productId"`
	Rating    int                `bson:"rating" json:"rating"`
	Comment   string             `bson:"comment" json:"comment"`
	CreatedAt int64              `bson:"created_at" json:"createdAt"`
}