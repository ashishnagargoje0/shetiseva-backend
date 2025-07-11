package controllers

import (
	"net/http"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var voiceFeedbackCollection *mongo.Collection

func InitVoiceFeedbackCollection() {
	voiceFeedbackCollection = config.DB.Collection("voice_feedback")
}

func SubmitVoiceFeedback(c *gin.Context) {
	var input models.VoiceFeedbackInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedback := bson.M{
		"userId":      input.UserID,
		"audioUrl":    input.AudioURL,
		"comment":     input.Comment,
		"submittedAt": time.Now(),
	}

	_, err := voiceFeedbackCollection.InsertOne(c, feedback)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save feedback"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voice feedback received"})
}
