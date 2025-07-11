package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/controllers"
)

func AdvisoryRoutes(r *gin.Engine) {
	r.GET("/advisory/crop/:cropId", controllers.GetCropAdvisory)
}
