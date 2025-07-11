package controllers

import (
	"net/http"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var reviewCollection *mongo.Collection

// InitReviewCollection initializes the MongoDB collection for reviews
func InitReviewCollection() {
	reviewCollection = config.DB.Collection("reviews")
}

func SubmitReview(c *gin.Context) {
	var input models.ReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, exists := c.Get("user_id") // ✅ Fix: match middleware key
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDVal.(primitive.ObjectID) // ✅ already ObjectID in middleware
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	review := bson.M{
		"userId":    userID,
		"productId": input.ProductID,
		"rating":    input.Rating,
		"comment":   input.Comment,
		"createdAt": time.Now().Unix(),
	}

	_, err := reviewCollection.InsertOne(c, review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review submitted successfully"})
}

func GetProductReviews(c *gin.Context) {
	productId := c.Param("id")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is required"})
		return
	}

	cursor, err := reviewCollection.Find(c, bson.M{"productId": productId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reviews"})
		return
	}

	var reviews []bson.M
	if err := cursor.All(c, &reviews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

