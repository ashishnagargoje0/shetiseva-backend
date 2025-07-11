package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
)

func ContactRoutes(router *gin.Engine) {
	// Public route
	router.POST("/contact", controllers.SubmitContactForm)

	// Admin protected
	admin := router.Group("/admin")
	admin.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		admin.GET("/contacts", controllers.GetAllContacts)
		admin.DELETE("/contacts/:id", controllers.DeleteContact)
	}
}
