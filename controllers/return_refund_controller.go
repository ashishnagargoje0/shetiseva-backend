package controllers

import (
	"net/http"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Declare collections as uninitialized
var returnRequestCollection *mongo.Collection
var refundRequestCollection *mongo.Collection

// Called from main.go AFTER DB is connected
func InitReturnRefundCollections() {
	returnRequestCollection = config.DB.Collection("return_requests")
	refundRequestCollection = config.DB.Collection("refund_requests")
}

// POST /order/return
func SubmitReturnRequest(c *gin.Context) {
	var input models.ReturnRequestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.OrderID == "" || input.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID and reason are required"})
		return
	}

	orderID, err := primitive.ObjectIDFromHex(input.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID format"})
		return
	}

	returnReq := models.ReturnRequest{
		ID:        primitive.NewObjectID(),
		OrderID:   orderID,
		Reason:    input.Reason,
		CreatedAt: time.Now().Unix(),
		Status:    "pending",
	}

	_, err = returnRequestCollection.InsertOne(c, returnReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit return request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Return request submitted", "returnRequestId": returnReq.ID.Hex()})
}

// POST /refund/request
func SubmitRefundRequest(c *gin.Context) {
	var input models.RefundRequestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.OrderID == "" || input.Amount <= 0 || input.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID, amount, and reason are required"})
		return
	}

	orderID, err := primitive.ObjectIDFromHex(input.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID format"})
		return
	}

	refundReq := models.RefundRequest{
		ID:        primitive.NewObjectID(),
		OrderID:   orderID,
		Amount:    input.Amount,
		Reason:    input.Reason,
		CreatedAt: time.Now().Unix(),
		Status:    "pending",
	}

	_, err = refundRequestCollection.InsertOne(c, refundReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit refund request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Refund request submitted", "refundRequestId": refundReq.ID.Hex()})
}
