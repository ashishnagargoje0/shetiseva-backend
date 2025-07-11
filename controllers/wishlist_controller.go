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

// Request body struct for AddToWishlist, with product_id as string from JSON
type AddWishlistRequest struct {
	ProductID string `json:"product_id" binding:"required"`
}

func AddToWishlist(c *gin.Context) {
	var req AddWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userObjID, ok := userIDValue.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	productObjID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	wishlistItem := models.WishlistItem{
		ID:        primitive.NewObjectID(),
		UserID:    userObjID,
		ProductID: productObjID,
	}

	collection := database.GetCollection("wishlist")

	existing := collection.FindOne(context.TODO(), bson.M{
		"user_id":    wishlistItem.UserID,
		"product_id": wishlistItem.ProductID,
	})
	if existing.Err() == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product already in wishlist"})
		return
	}

	_, err = collection.InsertOne(context.TODO(), wishlistItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to wishlist"})
}

func GetWishlist(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userObjID, ok := userIDValue.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	collection := database.GetCollection("wishlist")

	cursor, err := collection.Find(context.TODO(), bson.M{"user_id": userObjID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wishlist"})
		return
	}
	defer cursor.Close(context.TODO())

	var items []models.WishlistItem
	if err := cursor.All(context.TODO(), &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wishlist": items})
}

func RemoveFromWishlist(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection := database.GetCollection("wishlist")
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from wishlist"})
}
