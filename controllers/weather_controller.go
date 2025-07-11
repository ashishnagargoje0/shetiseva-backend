package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/ashishnagargoje0/backend/config"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTodayWeather(c *gin.Context) {
	district := c.Query("district")
	var weather models.WeatherData

	filter := bson.M{
		"date": time.Now().Format("2006-01-02"),
	}
	if district != "" {
		filter["district"] = district
	}

	err := config.DB.Collection("weather").FindOne(context.TODO(), filter).Decode(&weather)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No weather data for today"})
		return
	}

	c.JSON(http.StatusOK, weather)
}


func GetWeatherForecast(c *gin.Context) {
	district := c.Query("district")

	filter := bson.M{}
	if district != "" {
		filter["district"] = district
	}

	cursor, err := config.DB.Collection("weather_forecast").Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching forecast"})
		return
	}

	var forecast []models.WeatherForecast
	if err = cursor.All(context.TODO(), &forecast); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding forecast"})
		return
	}

	c.JSON(http.StatusOK, forecast)
}


