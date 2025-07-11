package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/controllers"
)

func AIRoutes(r *gin.Engine) {
	r.POST("/ai/diagnose", controllers.DiagnoseDisease)
	r.GET("/ai/alerts", controllers.GetAIAlerts)
}
