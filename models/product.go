package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	CategoryID  primitive.ObjectID `bson:"category_id" json:"category_id"`
	Price       float64            `bson:"price" json:"price"`
	ImageURL    string             `bson:"image_url" json:"image_url"`
	InStock     bool               `bson:"in_stock" json:"in_stock"`
	Tags        []string           `bson:"tags,omitempty" json:"tags,omitempty"`
}
