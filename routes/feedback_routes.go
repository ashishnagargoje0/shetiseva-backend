package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func FeedbackRoutes(r *gin.Engine) {
	feedback := r.Group("/feedback")
	feedback.Use(middlewares.AuthMiddleware())
	{
		feedback.POST("/voice", controllers.SubmitVoiceFeedback)           // Voice feedback endpoint
	}
}
