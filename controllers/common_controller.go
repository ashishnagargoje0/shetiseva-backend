package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 🏞 GET /district/list
func GetDistrictList(c *gin.Context) {
	districts := []string{
		"Ahmednagar", "Akola", "Amravati", "Aurangabad", "Beed", "Bhandara",
		"Buldhana", "Chandrapur", "Dhule", "Gadchiroli", "Gondia", "Hingoli",
		"Jalgaon", "Jalna", "Kolhapur", "Latur", "Mumbai", "Nagpur", "Nanded",
		"Nandurbar", "Nashik", "Osmanabad", "Palghar", "Parbhani", "Pune",
		"Raigad", "Ratnagiri", "Sangli", "Satara", "Sindhudurg", "Solapur",
		"Thane", "Wardha", "Washim", "Yavatmal",
	}

	c.JSON(http.StatusOK, gin.H{
		"districts": districts,
	})
}

// 🌐 Optional: GET /language/list
func GetSupportedLanguages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"languages": []string{"en", "hi", "mr"},
	})
}

// 💓 Optional: GET /health
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Shetiseva backend is running",
	})
}
