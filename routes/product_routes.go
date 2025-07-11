package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/controllers"
)

func ProductRoutes(r *gin.Engine) {
	r.GET("/products", controllers.GetAllProducts)
	r.GET("/products/filters", controllers.GetProductFilters)
	r.GET("/product/:id", controllers.GetProductByID)
	r.GET("/categories", controllers.GetAllCategories)
	r.GET("/category/:slug", controllers.GetCategoryProducts)
}
