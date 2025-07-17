package routes

import (
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func SubscriptionRoutes(router *gin.Engine) {
	// Group requiring user authentication
	authGroup := router.Group("/")
	authGroup.Use(middlewares.AuthMiddleware())
	{
		// Subscription routes
		authGroup.POST("/subscription/create", controllers.CreateSubscription)
		authGroup.POST("/subscription/pause", controllers.PauseSubscription)
		authGroup.POST("/subscription/resume", controllers.ResumeSubscription)
		authGroup.GET("/subscription/status", controllers.GetSubscriptionStatus)

		// Coins
		authGroup.GET("/coins/balance", controllers.GetCoinsBalance)

		// Referral
		authGroup.POST("/referral/use", controllers.UseReferralCode)
		authGroup.GET("/referral/stats", controllers.GetReferralStats)

		// Rewards
		authGroup.GET("/rewards/catalog", controllers.GetRewardsCatalog)
		authGroup.POST("/rewards/redeem", controllers.RedeemReward)
	}
}
