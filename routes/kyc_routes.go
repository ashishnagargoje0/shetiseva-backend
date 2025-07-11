package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func KYCRoutes(router *gin.Engine) {
	kyc := router.Group("/kyc")
	kyc.Use(middlewares.AuthMiddleware()) // 🔐 Protected KYC routes
	{
		kyc.POST("/upload", controllers.UploadKYC)    // Upload KYC document
		kyc.GET("/status", controllers.GetKYCStatus)  // Get current user's KYC status
		kyc.PUT("/verify/:email", middlewares.AdminMiddleware(), controllers.VerifyKYC)

	}
}
