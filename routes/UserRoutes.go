package routes

import (
	"github.com/ashishnagargoje0/shetiseva-backend/controllers"
	"github.com/ashishnagargoje0/shetiseva-backend/middleware"
	"github.com/gin-gonic/gin"
)

// UserRoutes defines all routes accessible to admin for user management
func UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")
	userGroup.Use(middleware.Auth(), middleware.AdminMiddleware()) // ✅ Use correct Auth() middleware

	{
		userGroup.GET("/", controllers.GetUsers)              // ✅ Fetch all users
		userGroup.GET("/:id", controllers.GetUserByID)        // ✅ Fetch user by ID
		userGroup.PUT("/:id/role", controllers.UpdateUserRole) // ✅ Update role
	}
}
