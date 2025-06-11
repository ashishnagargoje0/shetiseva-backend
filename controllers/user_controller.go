package controllers

import (
	"net/http"
	"strconv"

	"github.com/ashishnagargoje0/shetiseva-backend/config"
	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"github.com/ashishnagargoje0/shetiseva-backend/utils"
	"github.com/gin-gonic/gin"
)

// GET /users
func GetUsers(c *gin.Context) {
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	c.JSON(http.StatusOK, users)
}

// GET /users/:id
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		utils.HandleError(c, http.StatusNotFound, "User not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

// PUT /users/:id/role
func UpdateUserRole(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Role string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid role format")
		return
	}

	uid, err := strconv.Atoi(id)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := config.DB.Model(&models.User{}).Where("id = ?", uid).Update("role", input.Role).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to update role")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated"})
}

// ✅ NEW: GET /auth/profile (requires JWT)
func Profile(c *gin.Context) {
	userID := c.GetUint("user_id")
	email := c.GetString("email")

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   email,
	})
}
