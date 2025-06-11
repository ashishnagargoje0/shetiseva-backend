package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/ashishnagargoje0/shetiseva-backend/controllers"
    "github.com/ashishnagargoje0/shetiseva-backend/middleware"
)

func AuthRoutes(router *gin.Engine) {
    auth := router.Group("/auth")
    {
        auth.POST("/login", controllers.Login)
        auth.POST("/register", controllers.Register)
        auth.GET("/profile", middleware.Auth(), controllers.Profile) // Middleware + Protected route
    }
}
