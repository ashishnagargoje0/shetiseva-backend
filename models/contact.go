package models

import "time"

type Contact struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Message   string    `json:"message" bson:"message"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
