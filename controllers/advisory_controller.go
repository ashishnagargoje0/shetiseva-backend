package controllers

import (
    "context"
    "net/http"

    "github.com/ashishnagargoje0/backend/models"
    "github.com/ashishnagargoje0/backend/config"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
)

func GetCropAdvisory(c *gin.Context) {
    cropID := c.Param("cropId")
    var advisory models.CropAdvisory

    err := config.DB.Collection("advisories").FindOne(context.TODO(), bson.M{"crop_id": cropID}).Decode(&advisory)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Advisory not found"})
        return
    }

    c.JSON(http.StatusOK, advisory)
}
