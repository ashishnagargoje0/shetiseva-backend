package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/ashishnagargoje0/shetiseva-backend/config"
	"github.com/ashishnagargoje0/shetiseva-backend/controllers"
	"github.com/ashishnagargoje0/shetiseva-backend/middleware"
	"github.com/ashishnagargoje0/shetiseva-backend/routes"
)

func main() {
	// ✅ Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found. Proceeding with system/default config.")
	}

	// ✅ Set Gin to release mode if specified
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ✅ Connect to the database
	config.ConnectDB()
	log.Println("✅ Database connected successfully")

	// ✅ Initialize controllers with DB
	db := config.DB
	controllers.InitAuthController(db)
	controllers.InitProductController(db)
	controllers.InitCategoryController(db)
	controllers.InitWishlistController(db)
	controllers.InitContactController(db) // ✅ Don't forget this!

	// ✅ Create Gin router
	router := gin.Default()

	// ✅ Set trusted proxy
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("❌ Failed to set trusted proxies: %v", err)
	}

	// ✅ Serve static files (image uploads)
	router.Static("/uploads", "./uploads")

	// ✅ Apply middlewares
	router.Use(middleware.RateLimitMiddleware())
	// router.Use(middleware.CORSMiddleware())          // Uncomment if needed
	// router.Use(middleware.SecureHeadersMiddleware()) // Uncomment if needed

	// ✅ Register all routes from one place
	routes.SetupRoutes(router)

	// ✅ Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
