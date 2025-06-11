package routes

import (
	"github.com/ashishnagargoje0/shetiseva-backend/controllers"
	"github.com/ashishnagargoje0/shetiseva-backend/middleware"
	"github.com/gin-gonic/gin"
)

func CartRoutes(router *gin.Engine) {
	cart := router.Group("/cart").Use(middleware.Auth()) // ✅ FIXED
	{
		cart.GET("/", controllers.ViewCart)
		cart.POST("/", controllers.AddToCart)
		cart.PUT("/:id", controllers.UpdateCartItem)
		cart.DELETE("/:id", controllers.RemoveCartItem)
	}
}
