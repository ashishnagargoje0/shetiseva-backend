package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/ashishnagargoje0/backend/database"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ✅ Add this global variable
var userCollection *mongo.Collection

// ✅ Call this once in main.go after db.ConnectDB()
func InitUserCollection() {
	userCollection = database.UserCollection
}

// ✅ GetUserProfile
func GetUserProfile(c *gin.Context) {
	val, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userEmail, ok := val.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse email from token"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": userEmail}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Password = ""
	user.OTP = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// ✅ UpdateUserProfile
func UpdateUserProfile(c *gin.Context) {
	val, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userEmail := val.(string)

	var input struct {
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		ProfilePic string `json:"profile_pic"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	update := bson.M{}
	if input.Name != "" {
		update["name"] = input.Name
	}
	if input.Phone != "" {
		update["phone"] = input.Phone
	}
	if input.ProfilePic != "" {
		update["profile_pic"] = input.ProfilePic
	}

	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No update fields provided"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userCollection.UpdateOne(ctx, bson.M{"email": userEmail}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// ✅ SetUserLanguage
func SetUserLanguage(c *gin.Context) {
	val, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userEmail := val.(string)

	var input struct {
		Language string `json:"language" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language is required"})
		return
	}

	lang := strings.ToLower(input.Language)
	if lang != "en" && lang != "hi" && lang != "mr" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language code"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userCollection.UpdateOne(ctx,
		bson.M{"email": userEmail},
		bson.M{"$set": bson.M{"language": lang}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update language"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language updated", "language": lang})
}
