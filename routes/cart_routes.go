package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func CartRoutes(r *gin.Engine) {
    cartGroup := r.Group("/api/cart")
    cartGroup.Use(middlewares.AuthMiddleware()) // ✅ Add this
    {
        cartGroup.POST("/", controllers.AddToCart)
        cartGroup.GET("/", controllers.ViewCart)             // <-- Added this route for logged-in user
        cartGroup.GET("/:user_id", controllers.ViewCart)     // existing route (optional, or can be removed)
        cartGroup.DELETE("/", controllers.RemoveFromCart)
        cartGroup.PUT("/update-qty", controllers.UpdateCartQuantity)
    }

    // Checkout
    r.POST("/api/checkout", middlewares.AuthMiddleware(), controllers.Checkout)
    r.POST("/cart/add", controllers.AddToCart) // alias
    r.GET("/cart/view", controllers.ViewCartByQueryParam)
    r.DELETE("/cart/remove/:id", controllers.RemoveCartItemByID)
}


