package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func ReviewRoutes(r *gin.Engine) {
	// Public route — no auth
	r.GET("/review/product/:id", controllers.GetProductReviews)

	// Authenticated group
	review := r.Group("/review")
	review.Use(middlewares.AuthMiddleware())
	{
		review.POST("/submit", controllers.SubmitReview)
	}
}
