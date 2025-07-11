package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func SupportRoutes(r *gin.Engine) {
	support := r.Group("/support")
	support.Use(middlewares.AuthMiddleware())

	support.GET("/tickets", controllers.GetSupportTickets)           // ✅ View all support tickets
	support.POST("/ticket", controllers.SubmitSupportTicket)         // ✅ Submit a new ticket
	support.GET("/status/:id", controllers.GetSupportStatusByID)     // ✅ Get status by ticket ID
}
