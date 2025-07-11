package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// VoiceFeedback represents a user's voice feedback stored in the database.
type VoiceFeedback struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId" binding:"required"`
	AudioURL  string             `bson:"audio_url" json:"audioUrl" binding:"required,url"`
	CreatedAt int64              `bson:"created_at" json:"createdAt,omitempty"`
}
