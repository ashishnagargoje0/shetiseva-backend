package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"gorm.io/gorm"
)

var contactDB *gorm.DB

// InitContactController initializes the contact controller with the DB connection
func InitContactController(db *gorm.DB) {
	contactDB = db
}

// SubmitContactMessage handles POST /api/contact/
func SubmitContactMessage(c *gin.Context) {
	var input models.ContactMessage

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := contactDB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save contact message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact message submitted successfully"})
}

// GetContactMessages handles GET /admin/contacts
// This should be protected with AdminAuthMiddleware
func GetContactMessages(c *gin.Context) {
	var messages []models.ContactMessage

	if err := contactDB.Order("created_at desc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve contact messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
