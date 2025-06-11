package routes

import (
    "github.com/ashishnagargoje0/shetiseva-backend/controllers"
    "github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
    auth := r.Group("/auth")
    {
        auth.POST("/register", controllers.Register)
        auth.POST("/login", controllers.Login)
    }
}
