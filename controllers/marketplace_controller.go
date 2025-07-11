package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/database"
	"github.com/ashishnagargoje0/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SubmitMarketplaceAd(c *gin.Context) {
    var req models.MarketplaceAd

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
        return
    }

    // Check the type of userID and assign accordingly
    switch v := userID.(type) {
    case primitive.ObjectID:
        req.UserID = v
    case string:
        oid, err := primitive.ObjectIDFromHex(v)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }
        req.UserID = oid
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected user ID type"})
        return
    }

    req.ID = primitive.NewObjectID()

    collection := database.GetCollection("marketplace")

    _, err := collection.InsertOne(context.TODO(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit ad"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Ad submitted successfully"})
}


func GetAllMarketplaceItems(c *gin.Context) {
	collection := database.GetCollection("marketplace")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	defer cursor.Close(context.TODO())

	var items []models.MarketplaceAd
	if err := cursor.All(context.TODO(), &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func GetMarketplaceItemByID(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection := database.GetCollection("marketplace")

	var item models.MarketplaceAd
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&item)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}
