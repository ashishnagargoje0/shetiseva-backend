package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	db "github.com/ashishnagargoje0/backend/database"
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/internal/telemetry"
	"github.com/ashishnagargoje0/backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// ========== 1. Connect to MongoDB ==========
	config.ConnectMongoDB()
	if config.DB == nil {
		log.Fatal("❌ MongoDB not connected. Exiting...")
	}

	// ========== 2. Initialize Collections ==========
	db.ConnectDB() // If needed alongside config.ConnectMongoDB

	controllers.InitUserCollection()
	controllers.InitKYCCollection()
	controllers.InitCartCollection()
	controllers.InitOrderCollection()
	controllers.InitReturnRefundCollections()      // ✅ Added for return/refund
	controllers.InitSupportCollection()            // ✅ Added for support tickets
	controllers.InitVoiceFeedbackCollection()      // ✅ Added for voice feedback
	controllers.InitReviewCollection() 

	// ========== 3. Database Setup ==========
	db.InitDatabase()
	db.RunMigrations()
	db.SeedDatabase()

	// ========== 4. OpenTelemetry ==========
	shutdown := telemetry.InitTracer("shetiseva-backend")
	defer shutdown(context.Background())
	telemetry.InitMetrics()

	// ========== 5. Setup Gin ==========
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(telemetry.TracingMiddleware("shetiseva-backend"))
	router.Use(telemetry.MetricsMiddleware())
	telemetry.RegisterMetricsEndpoint(router)

	// Optional: Add CORS middleware here if frontend requires it

	// ========== 6. Health Check ==========
	router.GET("/health/db", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "MongoDB connected",
			"db":     config.DB.Name(),
		})
	})

	// ========== 7. Register All Routes ==========
	routes.AuthRoutes(router)
	routes.CartRoutes(router)
	routes.OrderRoutes(router)
	routes.ContactRoutes(router)
	routes.UserRoutes(router)
	routes.KYCRoutes(router)
	routes.CommonRoutes(router)
	routes.ProductRoutes(router)
	routes.WishlistRoutes(router)
	routes.CompareRoutes(router)
	routes.MarketplaceRoutes(router)
	routes.AdminRoutes(router)

	// ✅ NEW routes added for extended functionality
	routes.RefundRoutes(router)
	routes.SupportRoutes(router)
	routes.ReviewRoutes(router)
	routes.FeedbackRoutes(router)

	// ✅ NEW AI & advisory routes
	routes.AdvisoryRoutes(router)
	routes.WeatherRoutes(router)
	routes.AIRoutes(router)
	routes.ChatbotRoutes(router)

	// ========== 8. Start Server with Graceful Shutdown ==========
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()
	log.Println("🚀 Server started on http://localhost:8080")

	// ========== 9. Graceful Shutdown ==========
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🚦 Shutting down Shetiseva backend gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited properly")
}
