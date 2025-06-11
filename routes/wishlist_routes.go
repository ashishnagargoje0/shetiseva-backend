package routes

import (
	"github.com/ashishnagargoje0/shetiseva-backend/controllers"
	"github.com/ashishnagargoje0/shetiseva-backend/middleware"
	"github.com/gin-gonic/gin"
)

func WishlistRoutes(r *gin.Engine) {
	wishlist := r.Group("/wishlist")
	wishlist.Use(middleware.Auth())
	{
		wishlist.POST("/:product_id", controllers.AddToWishlist)
		wishlist.GET("/", controllers.GetWishlist)
		wishlist.DELETE("/:product_id", controllers.RemoveFromWishlist)
	}
}
