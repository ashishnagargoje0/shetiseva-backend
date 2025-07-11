package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ✅ Approve KYC (admin only)
func ApproveKYC(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"kyc_status": "approved"}}
	_, err := kycCollection.UpdateOne(ctx, bson.M{"email": body.Email}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve KYC"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KYC approved"})
}

// ❌ Reject KYC (admin only)
func RejectKYC(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"kyc_status": "rejected"}}
	_, err := kycCollection.UpdateOne(ctx, bson.M{"email": body.Email}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject KYC"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KYC rejected"})
}
