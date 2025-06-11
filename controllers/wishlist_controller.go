package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/ashishnagargoje0/shetiseva-backend/models"
)

var wishlistDB *gorm.DB

func InitWishlistController(db *gorm.DB) {
    wishlistDB = db
}

// ✅ Add to Wishlist
func AddToWishlist(c *gin.Context) {
    var input struct {
        ProductID uint `json:"product_id" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := c.GetUint("user_id")

    // Optional: prevent duplicates
    var existing models.Wishlist
    err := wishlistDB.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&existing).Error
    if err == nil {
        c.JSON(http.StatusConflict, gin.H{"message": "Product already in wishlist"})
        return
    }

    wishlist := models.Wishlist{
        UserID:    userID,
        ProductID: input.ProductID,
    }

    if err := wishlistDB.Create(&wishlist).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, wishlist)
}

// ✅ Get Wishlist
func GetWishlist(c *gin.Context) {
    userID := c.GetUint("user_id")
    var wishlist []models.Wishlist

    if err := wishlistDB.Preload("Product").Where("user_id = ?", userID).Find(&wishlist).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, wishlist)
}

// ✅ Remove from Wishlist
func RemoveFromWishlist(c *gin.Context) {
    userID := c.GetUint("user_id")
    productIDParam := c.Param("product_id")

    productID, err := strconv.Atoi(productIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id"})
        return
    }

    if err := wishlistDB.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.Wishlist{}).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Removed from wishlist"})
}
