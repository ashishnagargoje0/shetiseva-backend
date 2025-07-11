package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/ashishnagargoje0/backend/controllers"
    "github.com/ashishnagargoje0/backend/middlewares"
)

func MarketplaceRoutes(r *gin.Engine) {
    marketplace := r.Group("/marketplace")
    marketplace.Use(middlewares.AuthMiddleware())  // Optional: only if you want auth here

    marketplace.POST("/submit-ad", controllers.SubmitMarketplaceAd)
    marketplace.GET("/items", controllers.GetAllMarketplaceItems)
    marketplace.GET("/item/:id", controllers.GetMarketplaceItemByID)
}
