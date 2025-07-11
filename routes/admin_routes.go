package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	admin.Use(middlewares.AdminMiddleware()) // 🔐 Only admins allowed

	admin.POST("/kyc/approve", controllers.ApproveKYC)
	admin.POST("/kyc/reject", controllers.RejectKYC)
}
