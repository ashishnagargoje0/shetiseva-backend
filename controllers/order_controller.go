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
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection

// InitOrderCollection initializes the orders collection; call this once after DB connection is ready
func InitOrderCollection() {
	orderCollection = config.DB.Collection("orders")
}

// ======================= ORDER =======================

// OrderCheckout handles order placement from the user's cart.
func OrderCheckout(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDVal.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cartCol := config.DB.Collection("carts")
	cursor, err := cartCol.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}
	defer cursor.Close(ctx)

	var cartItems []models.CartItem
	if err := cursor.All(ctx, &cartItems); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode cart items"})
		return
	}
	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	order := models.Order{
		UserID:    userID,
		Items:     cartItems,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	_, err = orderCollection.InsertOne(ctx, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place order"})
		return
	}

	_, err = cartCol.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": order})
}

// GetOrderHistory fetches all orders for the logged-in user.
func GetOrderHistory(c *gin.Context) {
	userID, ok := getUserObjectID(c)
	if !ok {
		return // getUserObjectID already sends error
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := orderCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrderByID fetches a single order by its ID.
func GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var order models.Order
	err = orderCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ======================= PAYMENT =======================

// InitiatePayment simulates payment initiation (e.g., Razorpay).
func InitiatePayment(c *gin.Context) {
	var req struct {
		OrderID string `json:"orderId"`
		Method  string `json:"paymentGateway"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Dummy payment ID generation
	paymentID := primitive.NewObjectID().Hex()

	c.JSON(http.StatusOK, gin.H{
		"message":    "Payment initiated",
		"payment_id": paymentID,
		"gateway":    req.Method,
	})
}

// VerifyPayment updates payment status.
func VerifyPayment(c *gin.Context) {
	var req struct {
		OrderID   string `json:"orderId"`
		PaymentID string `json:"paymentId"`
		Status    string `json:"status"` // "success" / "failed"
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	orderID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if req.Status != "success" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment failed"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = orderCollection.UpdateOne(ctx, bson.M{"_id": orderID}, bson.M{"$set": bson.M{
		"status":         "paid",
		"payment_status": "success",
		"updated_at":     time.Now(),
	}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment verified successfully"})
}

// GetPaymentOptions returns available payment methods.
func GetPaymentOptions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"methods": []string{"COD", "UPI", "CreditCard", "EMI"},
	})
}

// ======================= EMI =======================

// ApplyEMIPlan simulates EMI application.
func ApplyEMIPlan(c *gin.Context) {
	var req struct {
		OrderID string `json:"orderId"`
		Months  int    `json:"months"`
	}
	if err := c.BindJSON(&req); err != nil || req.Months <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "EMI application submitted",
		"orderId":      req.OrderID,
		"months":       req.Months,
		"status":       "under_review",
		"interestRate": 12.5,
	})
}

// GetEMIStatus returns simulated EMI plan status.
func GetEMIStatus(c *gin.Context) {
	orderID := c.Query("orderId")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orderId":      orderID,
		"emiApproved":  true,
		"months":       6,
		"monthlyEMI":   1200,
		"status":       "approved",
		"interestRate": 12.5,
	})
}

// ======================= INVOICE =======================

// GetInvoice returns a simulated invoice URL.
func GetInvoice(c *gin.Context) {
	orderID := c.Param("orderId")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"invoice_url": "https://cdn.shetiseva.in/invoices/" + orderID + ".pdf",
	})
}

// ======================= DELIVERY =======================

// GetDeliveryAgentStatus returns agent's live location.
func GetDeliveryAgentStatus(c *gin.Context) {
	orderID := c.Query("orderId")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orderId":        orderID,
		"agentName":      "Ravi Pawar",
		"location":       "Nagpur",
		"eta":            "15 mins",
		"status":         "On the way",
		"contactNumber":  "+91-9876543210",
	})
}

// ConfirmDelivery marks the order as delivered.
func ConfirmDelivery(c *gin.Context) {
	var req struct {
		OrderID string `json:"orderId"`
	}
	if err := c.BindJSON(&req); err != nil || req.OrderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	orderID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = orderCollection.UpdateOne(ctx, bson.M{"_id": orderID}, bson.M{"$set": bson.M{
		"status":      "delivered",
		"deliveredAt": time.Now(),
	}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm delivery"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order marked as delivered"})
}
