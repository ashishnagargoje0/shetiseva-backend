package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")

	{
		// 🔐 Email + Password Auth
		auth.POST("/signup", controllers.Signup)              // Register with email
		auth.POST("/login", controllers.Login)                // Login with email/password
		auth.POST("/logout", controllers.Logout)              // Logout (invalidate token)

		// 🔑 Password Recovery
		auth.POST("/forgot-password", controllers.ForgotPassword)  // Send reset link/code
		auth.POST("/reset-password", controllers.ResetPassword)    // Reset password with code

		// 📱 Phone + OTP Auth
		auth.POST("/register-phone", controllers.RegisterPhone)    // Send OTP to phone
		auth.POST("/verify-otp", controllers.VerifyOTP)            // Verify OTP and login/signup
	}
}
