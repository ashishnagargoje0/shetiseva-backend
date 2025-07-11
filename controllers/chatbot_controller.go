package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/ashishnagargoje0/backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// POST /chatbot/message
func ChatWithBot(c *gin.Context) {
	var input models.ChatMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ✅ Get user ID from context (set in middleware)
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id missing from context"})
		return
	}
	userID, ok := userIDRaw.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Generate bot reply
	response := "🤖 Reply to: " + input.Message

	chat := models.ChatMessage{
		UserID:    userID,
		Message:   input.Message,
		Reply:     response,
		Timestamp: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := config.DB.Collection("chat_history").InsertOne(ctx, chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reply": response})
}

// GET /chatbot/history
func GetChatHistory(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, ok := userIDRaw.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.DB.Collection("chat_history").Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching chat history"})
		return
	}
	defer cursor.Close(ctx)

	var history []models.ChatMessage
	if err := cursor.All(ctx, &history); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding chat history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
