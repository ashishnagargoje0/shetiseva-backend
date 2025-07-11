package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func RefundRoutes(r *gin.Engine) {
	refund := r.Group("/")
	refund.Use(middlewares.AuthMiddleware())
	{
		refund.POST("/order/return", controllers.SubmitReturnRequest)     // Request product return
		refund.POST("/refund/request", controllers.SubmitRefundRequest)   // Request refund
	}
}
