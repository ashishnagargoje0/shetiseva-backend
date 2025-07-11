package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/database"
	"github.com/ashishnagargoje0/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ✅ POST /contact (Public)
func SubmitContactForm(c *gin.Context) {
	var input models.Contact
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	input.ID = primitive.NewObjectID().Hex()
	input.CreatedAt = time.Now()

	_, err := database.ContactCollection.InsertOne(context.TODO(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save contact message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message submitted successfully"})
}

// ✅ GET /admin/contacts (Admin only)
func GetAllContacts(c *gin.Context) {
	cursor, err := database.ContactCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}
	defer cursor.Close(context.TODO())

	var contacts []models.Contact
	if err := cursor.All(context.TODO(), &contacts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode messages"})
		return
	}

	c.JSON(http.StatusOK, contacts)
}

// ✅ DELETE /admin/contacts/:id (Admin only)
func DeleteContact(c *gin.Context) {
	id := c.Param("id")

	_, err := database.ContactCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}
