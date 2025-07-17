package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Subscription represents a user's subscription plan
type Subscription struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	PlanID     primitive.ObjectID `bson:"plan_id" json:"plan_id"`
	Status     string             `bson:"status" json:"status"` // active, paused, cancelled etc.
	StartedAt  time.Time          `bson:"started_at" json:"started_at"`
	PausedAt   *time.Time         `bson:"paused_at,omitempty" json:"paused_at,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	ModifiedAt time.Time          `bson:"modified_at" json:"modified_at"`
}

// SubscriptionInput used for creating a subscription
type SubscriptionInput struct {
	PlanID primitive.ObjectID `json:"plan_id" binding:"required"`
}

