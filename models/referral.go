package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// ReferralUseInput is the input for using a referral code
type ReferralUseInput struct {
    ReferralCode string `json:"referral_code" binding:"required"`
}

// Reward represents a redeemable reward item
type Reward struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name        string             `bson:"name" json:"name"`
    Description string             `bson:"description" json:"description"`
    Cost        int                `bson:"cost" json:"cost"` // cost in coins
    CreatedAt   primitive.DateTime `bson:"created_at,omitempty" json:"created_at,omitempty"`
    UpdatedAt   primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// RewardRedeemInput is input to redeem a reward
type RewardRedeemInput struct {
    RewardID primitive.ObjectID `json:"reward_id" binding:"required"`
}
