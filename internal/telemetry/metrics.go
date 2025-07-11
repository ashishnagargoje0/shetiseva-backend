package telemetry

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(RequestCount)
}

// Gin middleware to capture metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := http.StatusText(c.Writer.Status())

		RequestDuration.WithLabelValues(c.Request.Method, c.FullPath(), status).Observe(duration)
		RequestCount.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
	}
}

// Expose /metrics route
func RegisterMetricsEndpoint(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
