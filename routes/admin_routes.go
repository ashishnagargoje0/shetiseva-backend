package routes

import (
	"github.com/ashishnagargoje0/shetiseva-backend/controllers"
	"github.com/ashishnagargoje0/shetiseva-backend/middleware"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.Auth(), middleware.AdminOnly())

	// Dashboard stats
	adminGroup.GET("/dashboard-stats", controllers.GetAdminDashboardStats)

	// User management
	adminGroup.GET("/users", controllers.GetAllUsers)         // List all users
	adminGroup.PUT("/users/:id", controllers.UpdateUser)      // Update user info
	adminGroup.DELETE("/users/:id", controllers.DeleteUser)   // Delete user

	// Contact message management
	adminGroup.GET("/contacts", controllers.GetContactMessages)

	// Product management
	adminGroup.POST("/products", controllers.CreateProduct)
	adminGroup.GET("/products", controllers.GetProducts)
	adminGroup.GET("/products/low-stock", controllers.GetLowStockProducts)
	adminGroup.GET("/products/:id", controllers.GetProduct)
	adminGroup.DELETE("/products/:id", controllers.DeleteProduct)

	// Category management
	adminGroup.POST("/categories", controllers.CreateCategory)
	adminGroup.GET("/categories", controllers.GetCategories)
	adminGroup.GET("/categories/:id", controllers.GetCategoryByID)
	adminGroup.PUT("/categories/:id", controllers.UpdateCategory)
	adminGroup.DELETE("/categories/:id", controllers.DeleteCategory)
}
