package controllers

import (
	"net/http"

	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"github.com/ashishnagargoje0/shetiseva-backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var authDB *gorm.DB

// InitAuthController initializes the auth controller with DB
func InitAuthController(db *gorm.DB) {
	authDB = db
}

// Register handles user signup
func Register(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Check for existing user
	var existing models.User
	if err := authDB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		utils.HandleError(c, http.StatusConflict, "Email already registered")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Error hashing password")
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := authDB.Create(&user).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handles user login and token generation
func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid input")
		return
	}

	var user models.User
	if err := authDB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Incorrect password")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Token generation failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
