package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(r *gin.Engine) {
	orders := r.Group("/api/orders")
	orders.Use(middlewares.AuthMiddleware())
	{
		// Order endpoints
		orders.POST("/", controllers.OrderCheckout)         // Place new order
		orders.GET("/", controllers.GetOrderHistory)        // User's order history
		orders.GET("/:id", controllers.GetOrderByID)        // Get specific order by ID

		// Payment endpoints
		orders.POST("/payment/initiate", controllers.InitiatePayment)  // Start payment
		orders.POST("/payment/verify", controllers.VerifyPayment)      // Confirm payment
		orders.GET("/payment/options", controllers.GetPaymentOptions)  // Available payment methods

		// EMI endpoints
		orders.POST("/emiplan/apply", controllers.ApplyEMIPlan)        // Apply for EMI
		orders.GET("/emiplan/status", controllers.GetEMIStatus)        // EMI application status

		// Invoice endpoint
		orders.GET("/invoice/:orderId", controllers.GetInvoice)        // Get invoice PDF link

		// Delivery endpoints
		orders.GET("/delivery/agent-status", controllers.GetDeliveryAgentStatus) // Track agent
		orders.POST("/delivery/confirm", controllers.ConfirmDelivery)            // Confirm delivery
	}
}
