package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/shetiseva-backend/controllers"
)

func ContactRoutes(router *gin.Engine) {
	api := router.Group("/api/contact")
	{
		api.POST("/", controllers.SubmitContactMessage)
	}
}
