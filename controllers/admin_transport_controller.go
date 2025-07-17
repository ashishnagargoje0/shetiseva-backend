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

func UpdateTransportStatus(c *gin.Context) {
	var input models.UpdateTransportStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	bookingID, err := primitive.ObjectIDFromHex(input.BookingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid booking ID"})
		return
	}

	collection := config.DB.Collection("transport_bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"status": input.Status}}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": bookingID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Booking not found or update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transport status updated successfully"})
}
