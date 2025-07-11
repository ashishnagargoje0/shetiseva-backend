package controllers

import (
	"net/http"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var supportCollection *mongo.Collection

func InitSupportCollection() {
	supportCollection = config.DB.Collection("support_tickets")
}
// POST /support/ticket
func SubmitSupportTicket(c *gin.Context) {
	var input models.SupportTicketInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, exists := c.Get("user_id") // ✅ match key set in middleware
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDVal.(primitive.ObjectID) // ✅ already ObjectID in middleware
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	ticket := bson.M{
		"userId":    userID,
		"subject":   input.Subject,
		"message":   input.Message,
		"status":    "open",
		"createdAt": time.Now().Unix(),
	}

	_, err := supportCollection.InsertOne(c, ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Support ticket submitted successfully"})
}



// GET /support/tickets
func GetSupportTickets(c *gin.Context) {
	cursor, err := supportCollection.Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tickets"})
		return
	}
	var tickets []bson.M
	if err = cursor.All(c, &tickets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing tickets"})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

// GET /support/status/:id
func GetSupportStatusByID(c *gin.Context) {
	idParam := c.Param("id")
	ticketID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}
	var ticket bson.M
	err = supportCollection.FindOne(c, bson.M{"_id": ticketID}).Decode(&ticket)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}
	c.JSON(http.StatusOK, ticket)
}
