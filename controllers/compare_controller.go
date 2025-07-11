package controllers

import (
	"context"
	"net/http"

	"github.com/ashishnagargoje0/backend/database"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddToCompare(c *gin.Context) {
	var req models.CompareItem

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 🛠 Safe conversion
	var uid primitive.ObjectID
	switch id := userID.(type) {
	case string:
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		uid = objID
	case primitive.ObjectID:
		uid = id
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected user ID type"})
		return
	}

	req.UserID = uid
	req.ID = primitive.NewObjectID()

	collection := database.GetCollection("compare")

	// 🔍 Check for duplicates
	existsDoc := collection.FindOne(context.TODO(), bson.M{
		"user_id":    req.UserID,
		"product_id": req.ProductID,
	})
	if existsDoc.Err() == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product already in compare list"})
		return
	}

	// ✅ Insert
	_, err := collection.InsertOne(context.TODO(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to compare list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to compare list"})
}


func GetCompareList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var uid primitive.ObjectID
	switch id := userID.(type) {
	case string:
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		uid = objID
	case primitive.ObjectID:
		uid = id
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected user ID type"})
		return
	}

	collection := database.GetCollection("compare")

	cursor, err := collection.Find(context.TODO(), bson.M{"user_id": uid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch compare list"})
		return
	}
	defer cursor.Close(context.TODO())

	var items []models.CompareItem
	if err := cursor.All(context.TODO(), &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse compare list"})
		return
	}

	c.JSON(http.StatusOK, items)
}
