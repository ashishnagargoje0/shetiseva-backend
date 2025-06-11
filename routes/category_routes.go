package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/ashishnagargoje0/shetiseva-backend/controllers"
)

func CategoryRoutes(router *gin.Engine) {
    router.POST("/categories/", controllers.CreateCategory)
    router.GET("/categories/", controllers.GetCategories)
    router.GET("/categories/:id", controllers.GetCategoryByID)
    router.PUT("/categories/:id", controllers.UpdateCategory)
    router.DELETE("/categories/:id", controllers.DeleteCategory)
}
