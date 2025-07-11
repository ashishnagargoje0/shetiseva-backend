package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
)

func WishlistRoutes(r *gin.Engine) {
	wishlist := r.Group("/wishlist").Use(middlewares.AuthMiddleware())
	{
		wishlist.POST("/add", controllers.AddToWishlist)
		wishlist.GET("/", controllers.GetWishlist)
		wishlist.DELETE("/:id", controllers.RemoveFromWishlist)
	}
}
