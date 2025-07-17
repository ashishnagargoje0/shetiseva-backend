package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/controllers"
)

func MandiRoutes(r *gin.Engine) {
	// /mandi routes
	mandi := r.Group("/mandi")
	{
		mandi.GET("/prices", controllers.GetMandiPrices)
		mandi.POST("/alerts", controllers.SetPriceAlert)
	}

	// /transport routes
	transport := r.Group("/transport")
	{
		transport.POST("/book", controllers.BookTransport)
		transport.GET("/status/:id", controllers.GetTransportStatus)
		r.POST("/mandi/alerts/cancel", controllers.CancelPriceAlert)
	}
}
