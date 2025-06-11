package routes

import (
	"github.com/ashishnagargoje0/shetiseva-backend/controllers"
	"github.com/ashishnagargoje0/shetiseva-backend/middleware"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.Engine) {
	// Authenticated user routes
	order := router.Group("/orders").Use(middleware.Auth())
	{
		order.POST("/", controllers.CreateOrder)
		order.GET("/", controllers.GetAllOrders)     // ✅ List user orders
		order.GET("/:id", controllers.GetOrder)
		order.POST("/:id/pay", controllers.PayOrder) // ✅ Step 2 - Payment placeholder
	}

	// Admin routes (you can protect with admin middleware later)
	admin := router.Group("/admin/orders")
	{
		admin.GET("/", controllers.AdminGetAllOrders)              // ✅ View all orders
		admin.PUT("/:id/status", controllers.AdminUpdateOrderStatus) // ✅ Change order status
	}
}
