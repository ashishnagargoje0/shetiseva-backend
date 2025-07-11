package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	user := router.Group("/user")
	user.Use(middlewares.AuthMiddleware()) // 🔒 Protected routes
	{
		user.GET("/me", controllers.GetUserProfile)
		user.POST("/profile/update", controllers.UpdateUserProfile)
		user.POST("/language/set", controllers.SetUserLanguage)
	}
}
// This sets up the user routes with authentication middleware.