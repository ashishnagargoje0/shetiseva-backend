package routes

import (
	"github.com/ashishnagargoje0/shetiseva-backend/controllers"
	"github.com/ashishnagargoje0/shetiseva-backend/middleware"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.Engine) {
	// Public product routes
	public := r.Group("/products")
	{
		public.GET("/", controllers.GetProducts)
		public.GET("/:id", controllers.GetProduct)
	}

	// Protected product routes (require JWT auth)
	protected := r.Group("/products")
	protected.Use(middleware.Auth())
	{
		protected.POST("/", controllers.CreateProduct)
		protected.PUT("/:id", controllers.UpdateProduct)
		protected.DELETE("/:id", controllers.DeleteProduct)
		protected.GET("/low-stock", controllers.GetLowStockProducts)
	}
}
