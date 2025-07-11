package telemetry

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log after response
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("[GIN] %v | %3d | %13v | %s | %s",
			start.Format("2006/01/02 - 15:04:05"),
			status,
			latency,
			method,
			path,
		)
	}
}
