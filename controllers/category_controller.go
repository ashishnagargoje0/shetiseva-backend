package controllers

import (
	"net/http"

	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var categoryDB *gorm.DB

// InitCategoryController initializes DB instance for category operations
func InitCategoryController(db *gorm.DB) {
	categoryDB = db
}

// CreateCategory handles POST /categories/
func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := categoryDB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategories handles GET /categories/
func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := categoryDB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryByID handles GET /categories/:id
func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := categoryDB.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category"})
		}
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory handles PUT /categories/:id
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := categoryDB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var updatedData models.Category
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.Name = updatedData.Name
	category.Description = updatedData.Description

	if err := categoryDB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory handles DELETE /categories/:id
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := categoryDB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
