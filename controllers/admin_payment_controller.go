package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/ashishnagargoje0/backend/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GET /admin/invoices/:id
func GetAdminInvoiceByID(c *gin.Context) {
	invoiceID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(invoiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	var invoice models.Invoice
	err = database.InvoiceCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&invoice)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}


// POST /admin/refund/initiate
func InitiateRefund(c *gin.Context) {
	var input struct {
		OrderID string  `json:"order_id" binding:"required"`
		Reason  string  `json:"reason" binding:"required"`
		Amount  float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refund := models.Refund{
		ID:        primitive.NewObjectID(),
		OrderID:   input.OrderID,
		Reason:    input.Reason,
		Amount:    input.Amount,
		Status:    "Initiated",
		CreatedAt: time.Now(),
	}

	_, err := database.RefundCollection.InsertOne(context.TODO(), refund)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate refund"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Refund initiated", "refund_id": refund.ID.Hex()})
}

// POST /admin/commission/calculate
func CalculateCommission(c *gin.Context) {
	var input struct {
		OrderID    string  `json:"order_id" binding:"required"`
		Percentage float64 `json:"percentage" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Convert string to ObjectID
	orderObjectID, err := primitive.ObjectIDFromHex(input.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	err = database.OrderCollection.FindOne(context.TODO(), bson.M{"_id": orderObjectID}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	commissionAmount := (order.TotalAmount * input.Percentage) / 100

	c.JSON(http.StatusOK, gin.H{
		"order_id":          input.OrderID,
		"total_amount":      order.TotalAmount,
		"commission_rate":   input.Percentage,
		"commission_amount": commissionAmount,
	})
}


// GET /admin/payment/history
func GetAdminPaymentHistory(c *gin.Context) {
	cursor, err := database.PaymentCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment history"})
		return
	}
	defer cursor.Close(context.TODO())

	var payments []models.Payment
	if err := cursor.All(context.TODO(), &payments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payments": payments})
}

// POST /admin/credit/manual-add
func ManualCredit(c *gin.Context) {
	var input struct {
		UserID string  `json:"user_id" binding:"required"`
		Amount float64 `json:"amount" binding:"required"`
		Reason string  `json:"reason"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	credit := models.WalletTransaction{
		ID:        primitive.NewObjectID(),
		UserID:    input.UserID,
		Amount:    input.Amount,
		Type:      "Credit",
		Reason:    input.Reason,
		CreatedAt: time.Now(),
	}

	_, err := database.WalletTransactionCollection.InsertOne(context.TODO(), credit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add credit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manual credit added", "transaction_id": credit.ID.Hex()})
}

