package controllers

import (
    "net/http"
    "strconv"

    "github.com/ashishnagargoje0/shetiseva-backend/config"
    "github.com/ashishnagargoje0/shetiseva-backend/models"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)


// AddToCart adds a product to the user's cart or increases quantity if already present
func AddToCart(c *gin.Context) {
    userIDRaw, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    userID := userIDRaw.(uint)

    var input struct {
        ProductID uint `json:"product_id" binding:"required"`
        Quantity  uint `json:"quantity" binding:"required,min=1"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var cartItem models.CartItem
    err := config.DB.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&cartItem).Error
    if err != nil && err != gorm.ErrRecordNotFound {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    if err == gorm.ErrRecordNotFound {
        cartItem = models.CartItem{
            UserID:    userID,
            ProductID: input.ProductID,
            Quantity:  input.Quantity,
        }
        if err := config.DB.Create(&cartItem).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
            return
        }
    } else {
        cartItem.Quantity += input.Quantity
        if err := config.DB.Save(&cartItem).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{"message": "Added to cart", "cart_item": cartItem})
}

// UpdateCartItem updates the quantity of a cart item by ID
func UpdateCartItem(c *gin.Context) {
    userIDRaw, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    userID := userIDRaw.(uint)

    cartItemID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
        return
    }

    var input struct {
        Quantity uint `json:"quantity" binding:"required,min=1"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var cartItem models.CartItem
    if err := config.DB.First(&cartItem, uint(cartItemID)).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
        return
    }

    if cartItem.UserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this cart item"})
        return
    }

    cartItem.Quantity = input.Quantity
    if err := config.DB.Save(&cartItem).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Cart item updated", "cart_item": cartItem})
}

// RemoveCartItem deletes a cart item by ID
func RemoveCartItem(c *gin.Context) {
    userIDRaw, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    userID := userIDRaw.(uint)

    cartItemID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
        return
    }

    var cartItem models.CartItem
    if err := config.DB.First(&cartItem, uint(cartItemID)).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
        return
    }

    if cartItem.UserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this cart item"})
        return
    }

    if err := config.DB.Delete(&cartItem).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cart item"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Cart item removed"})
}

// ViewCart returns all cart items for the logged-in user
func ViewCart(c *gin.Context) {
    userIDRaw, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    userID := userIDRaw.(uint)

    var cartItems []models.CartItem
    if err := config.DB.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart items"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"cart": cartItems})
}
