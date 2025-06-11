package controllers

import (
	"net/http"
	"strconv"

	"github.com/ashishnagargoje0/shetiseva-backend/config"
	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"github.com/gin-gonic/gin"
)

// getUserID extracts the user_id from context
func getUserID(c *gin.Context) (uint, bool) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return 0, false
	}
	return userIDRaw.(uint), true
}

// CreateOrder creates a new order and connects to cart
func CreateOrder(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	// Fetch user's cart
	var cartItems []models.CartItem
	if err := config.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil || len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty or fetch failed"})
		return
	}

	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += float64(item.Quantity) * item.Product.Price
	}

	// Create order
	order := models.Order{
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      "pending",
	}
	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Create order items
	for _, item := range cartItems {
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
		}
		config.DB.Create(&orderItem)
	}

	// Clear cart
	config.DB.Where("user_id = ?", userID).Delete(&models.CartItem{})

	c.JSON(http.StatusCreated, gin.H{"message": "Order created", "order": order})
}

// GetOrder retrieves a specific order by ID
func GetOrder(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	if err := config.DB.Preload("OrderItems.Product").First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// GetAllOrders returns all orders of the logged-in user
func GetAllOrders(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var orders []models.Order
	if err := config.DB.Preload("OrderItems.Product").Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// PayOrder is a placeholder for payment processing
func PayOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = "paid"
	config.DB.Save(&order)

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful (placeholder)", "order": order})
}

// AdminGetAllOrders - View all orders
func AdminGetAllOrders(c *gin.Context) {
	var orders []models.Order
	if err := config.DB.Preload("User").Preload("OrderItems.Product").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all orders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// AdminUpdateOrderStatus - Change order status (admin)
func AdminUpdateOrderStatus(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = input.Status
	config.DB.Save(&order)

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated", "order": order})
}
