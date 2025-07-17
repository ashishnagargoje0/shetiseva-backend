package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/ashishnagargoje0/backend/database"
	"github.com/ashishnagargoje0/backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GET /mandi/prices
func GetMandiPrices(c *gin.Context) {
	crop := c.Query("crop")
	district := c.Query("district")

	filter := bson.M{}
	if crop != "" {
		filter["crop"] = crop
	}
	if district != "" {
		filter["district"] = district
	}

	cursor, err := database.MandiPriceCollection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prices"})
		return
	}
	defer cursor.Close(context.TODO())

	var prices []bson.M
	if err := cursor.All(context.TODO(), &prices); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding prices"})
		return
	}

	c.JSON(http.StatusOK, prices)
}

// POST /mandi/alerts
func SetPriceAlert(c *gin.Context) {
	var alert models.PriceAlertInput
	if err := c.ShouldBindJSON(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	alertDoc := bson.M{
		"user_id":         alert.UserID,
		"crop":            alert.Crop,
		"district":        alert.District,
		"price_threshold": alert.PriceThreshold,
		"notify_method":   alert.NotifyMethod,
		"created_at":      time.Now(),
	}

	_, err := database.MandiAlertCollection.InsertOne(context.TODO(), alertDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set alert"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alert set successfully"})
}

// POST /transport/book
func BookTransport(c *gin.Context) {
	var booking models.TransportBookingInput
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	bookingDoc := bson.M{
		"user_id":         booking.UserID,
		"crop":            booking.Crop,
		"quantity":        booking.Quantity,
		"pickup_location": booking.PickupLocation,
		"drop_location":   booking.DropLocation,
		"date":            booking.Date,
		"status":          "pending",
		"created_at":      time.Now(),
	}

	res, err := database.TransportBookingCollection.InsertOne(context.TODO(), bookingDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book transport"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking successful", "booking_id": res.InsertedID})
}

// GET /transport/status/:id
func GetTransportStatus(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var result bson.M
	err = database.TransportBookingCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}
func CancelPriceAlert(c *gin.Context) {
	var input models.CancelPriceAlertInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	collection := config.DB.Collection("mandi_alerts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":  input.UserID,
		"crop":     input.Crop,
		"district": input.District,
	}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil || res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No alert found to cancel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Price alert cancelled successfully"})
}
