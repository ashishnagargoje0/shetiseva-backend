package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ====== POST /subscription/create ======
func CreateSubscription(c *gin.Context) {
	var input models.SubscriptionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(primitive.ObjectID)

	subscription := models.Subscription{
		UserID:     userID,
		PlanID:     input.PlanID,
		Status:     "active",
		StartedAt:  time.Now(),
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	subscriptionCollection := config.DB.Collection("subscriptions")

	res, err := subscriptionCollection.InsertOne(ctx, subscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription created", "subscription_id": res.InsertedID})
}

// ====== POST /subscription/pause ======
func PauseSubscription(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(primitive.ObjectID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	subscriptionCollection := config.DB.Collection("subscriptions")

	filter := bson.M{"user_id": userID, "status": "active"}
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":      "paused",
			"paused_at":   now,
			"modified_at": now,
		},
	}

	result, err := subscriptionCollection.UpdateOne(ctx, filter, update)
	if err != nil || result.ModifiedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active subscription found or failed to pause"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription paused"})
}

// ====== POST /subscription/resume ======
func ResumeSubscription(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(primitive.ObjectID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	subscriptionCollection := config.DB.Collection("subscriptions")

	filter := bson.M{"user_id": userID, "status": "paused"}
	update := bson.M{
		"$set": bson.M{
			"status":      "active",
			"paused_at":   nil,
			"modified_at": time.Now(),
		},
	}

	result, err := subscriptionCollection.UpdateOne(ctx, filter, update)
	if err != nil || result.ModifiedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No paused subscription found or failed to resume"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription resumed"})
}

// ====== GET /subscription/status ======
func GetSubscriptionStatus(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(primitive.ObjectID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	subscriptionCollection := config.DB.Collection("subscriptions")

	var subscription models.Subscription
	err := subscriptionCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&subscription)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No subscription found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// ====== GET /coins/balance ======
func GetCoinsBalance(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(primitive.ObjectID)

	coinsCollection := config.DB.Collection("coins")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var coins models.CoinsBalance
	err := coinsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&coins)
	if err != nil {
		// If no record found, respond zero balance
		coins = models.CoinsBalance{
			UserID: userID,
			Coins:  0,
		}
	}

	c.JSON(http.StatusOK, gin.H{"coins": coins.Coins})
}

// ====== POST /referral/use ======
func UseReferralCode(c *gin.Context) {
	var input models.ReferralUseInput
	if err := c.ShouldBindJSON(&input); err != nil || input.ReferralCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid referral code"})
		return
	}

	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(primitive.ObjectID)

	referralCollection := config.DB.Collection("referrals")
	referralUsesCollection := config.DB.Collection("referral_uses")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if referral code exists and is valid
	var ref bson.M
	err := referralCollection.FindOne(ctx, bson.M{"code": input.ReferralCode}).Decode(&ref)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid referral code"})
		return
	}

	// Check if user already used this referral code
	count, _ := referralUsesCollection.CountDocuments(ctx, bson.M{"user_id": userID, "referral_code": input.ReferralCode})
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Referral code already used"})
		return
	}

	// Mark referral code used by user
	_, err = referralUsesCollection.InsertOne(ctx, bson.M{
		"user_id":       userID,
		"referral_code": input.ReferralCode,
		"referrer_id":   ref["user_id"],
		"used_at":       time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to use referral code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Referral code applied"})
}

// ====== GET /referral/stats ======
func GetReferralStats(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(primitive.ObjectID)

	referralUsesCollection := config.DB.Collection("referral_uses")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := referralUsesCollection.CountDocuments(ctx, bson.M{"referrer_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referral stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"referral_count": count})
}

// ====== GET /rewards/catalog ======
func GetRewardsCatalog(c *gin.Context) {
	rewardsCollection := config.DB.Collection("rewards")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := rewardsCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rewards"})
		return
	}
	defer cursor.Close(ctx)

	var rewards []models.Reward
	if err := cursor.All(ctx, &rewards); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode rewards"})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

// ====== POST /rewards/redeem ======
func RedeemReward(c *gin.Context) {
	var input models.RewardRedeemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid redeem request"})
		return
	}

	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(primitive.ObjectID)

	coinsCollection := config.DB.Collection("coins")
	rewardsCollection := config.DB.Collection("rewards")
	redemptionsCollection := config.DB.Collection("reward_redemptions")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find reward by ID
	var reward models.Reward
	err := rewardsCollection.FindOne(ctx, bson.M{"_id": input.RewardID}).Decode(&reward)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reward not found"})
		return
	}

	// Find user's coin balance
	var coins models.CoinsBalance
	err = coinsCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&coins)
	if err != nil || coins.Coins < reward.Cost {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient coins"})
		return
	}

	// Deduct coins
	update := bson.M{"$inc": bson.M{"coins": -reward.Cost}}
	_, err = coinsCollection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deduct coins"})
		return
	}

	// Record redemption
	_, err = redemptionsCollection.InsertOne(ctx, bson.M{
		"user_id":     userID,
		"reward_id":   reward.ID,
		"redeemed_at": time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record redemption"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reward redeemed successfully"})
}
