package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func CompareRoutes(r *gin.Engine) {
	group := r.Group("/compare").Use(middlewares.AuthMiddleware())
	{
		group.POST("/add", controllers.AddToCompare)
		group.GET("/view", controllers.GetCompareList)
	}
}
