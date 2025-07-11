package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/controllers"
)

func ChatbotRoutes(r *gin.Engine) {
	r.POST("/chatbot/message", controllers.ChatWithBot)
	r.GET("/chatbot/history", controllers.GetChatHistory)
}
