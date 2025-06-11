package controllers

import (
	"net/http"
	"strconv"

	"github.com/ashishnagargoje0/shetiseva-backend/config"
	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"github.com/ashishnagargoje0/shetiseva-backend/utils"
	"github.com/gin-gonic/gin"
)

// GetAllUsers returns all registered users (Admin only)
func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser allows admin to update user details
func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		utils.HandleError(c, http.StatusNotFound, "User not found")
		return
	}

	var input struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		IsAdmin *bool  `json:"is_admin"` // pointer to detect if field is provided
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.IsAdmin != nil {
		user.IsAdmin = *input.IsAdmin
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser removes a user by ID
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		utils.HandleError(c, http.StatusNotFound, "User not found")
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
