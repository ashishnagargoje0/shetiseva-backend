package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/gin-gonic/gin"
)

func CommonRoutes(router *gin.Engine) {
	common := router.Group("/")
	{
		common.GET("/district/list", controllers.GetDistrictList)
		common.GET("/language/list", controllers.GetSupportedLanguages) // optional
		common.GET("/health", controllers.HealthCheck)                   // optional
	}
}
