package controllers

import (
	"net/http"

	"github.com/ashishnagargoje0/shetiseva-backend/config"
	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"github.com/gin-gonic/gin"
)

func GetAdminDashboardStats(c *gin.Context) {
	var totalUsers int64
	var totalProducts int64
	var lowStockCount int64
	var totalCategories int64

	config.DB.Model(&models.User{}).Count(&totalUsers)
	config.DB.Model(&models.Product{}).Count(&totalProducts)
	config.DB.Model(&models.Product{}).Where("quantity <= ?", 10).Count(&lowStockCount)
	config.DB.Model(&models.Category{}).Count(&totalCategories)

	c.JSON(http.StatusOK, gin.H{
		"total_users":      totalUsers,
		"total_products":   totalProducts,
		"low_stock_items":  lowStockCount,
		"total_categories": totalCategories,
	})
}
