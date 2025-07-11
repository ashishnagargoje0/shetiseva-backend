package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ashishnagargoje0/backend/controllers"
)

func WeatherRoutes(r *gin.Engine) {
	r.GET("/weather/today", controllers.GetTodayWeather)
	r.GET("/weather/forecast", controllers.GetWeatherForecast)
}
