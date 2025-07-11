package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/ashishnagargoje0/backend/config"
	"go.mongodb.org/mongo-driver/bson"
)

func DiagnoseDisease(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}

	imagePath := "./temp/" + file.Filename
	_ = c.SaveUploadedFile(file, imagePath) // Save image locally (optional)

	// Simulated AI result (replace with real AI model integration later)
	result := "Leaf Blight Detected"

	c.JSON(http.StatusOK, gin.H{
		"diagnosis": result,
		"image":     file.Filename,
	})
}

func GetAIAlerts(c *gin.Context) {
	cursor, err := config.DB.Collection("ai_alerts").Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch alerts"})
		return
	}

	var alerts []models.AIAlert
	if err := cursor.All(context.TODO(), &alerts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding alerts"})
		return
	}

	c.JSON(http.StatusOK, alerts)
}
