package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func ConsultationRoutes(router *gin.Engine) {
	// 📞 Farmer Consultation Routes
	consultation := router.Group("/consultation").Use(middlewares.AuthMiddleware())
	{
		consultation.POST("/book", controllers.BookConsultation)
		consultation.GET("/status/:id", controllers.GetConsultationStatus)
		consultation.POST("/feedback", controllers.SubmitConsultationFeedback)
	}

	// 🚁 Drone Service Routes (User)
	drone := router.Group("/drone").Use(middlewares.AuthMiddleware())
	{
		drone.POST("/book", controllers.BookDrone)
		drone.GET("/availability", controllers.GetDroneAvailability)
		drone.POST("/cancel", controllers.CancelDroneBooking)
		drone.GET("/status", controllers.GetMyDroneBookings) // ✅ View personal bookings
	}

	// 🛠️ Admin-only Drone Management Routes
	adminDrone := router.Group("/admin/drone").Use(middlewares.AdminMiddleware())
	{
		adminDrone.GET("/bookings", controllers.AdminGetAllDroneBookings)    // ✅ View all
		adminDrone.PUT("/approve/:id", controllers.AdminApproveDroneBooking) // ✅ Approve booking
		adminDrone.PUT("/reject/:id", controllers.AdminRejectDroneBooking)   // ✅ Reject booking
	}

	// ✴️ Optional: Lightweight Admin Approval via POST
	droneAdmin := router.Group("/drone").Use(middlewares.AdminMiddleware())
	{
		droneAdmin.POST("/approve/:id", controllers.ApproveDroneBooking)
		droneAdmin.POST("/reject/:id", controllers.RejectDroneBooking)
	}
}
