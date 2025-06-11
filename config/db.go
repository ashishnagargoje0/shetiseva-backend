package config

import (
	"os"

	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"github.com/ashishnagargoje0/shetiseva-backend/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB     // Main production database
var TestDB *gorm.DB // Optional: used for automated tests

// ✅ ConnectDB sets up the main PostgreSQL connection and auto-migrates models
func ConnectDB() *gorm.DB {
	utils.InitLogger() // Initialize logger before DB actions

	// ✅ Get DATABASE_URL from environment or fallback to local default
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=1016 dbname=shetiseva port=5432 sslmode=disable"
		utils.InfoLogger.Println("⚠️ Environment variable DATABASE_URL not set. Using fallback DSN.")
	}

	// ✅ Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.ErrorLogger.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// ✅ Set global DB reference
	DB = db
	utils.InfoLogger.Println("✅ Successfully connected to the production database.")

	// ✅ Auto-migrate all models
	err = db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.CartItem{},
		&models.Wishlist{},         // ✅ Correct name
		&models.Order{},
		&models.OrderItem{},  
		&models.ContactMessage{},
	)
	if err != nil {
		utils.ErrorLogger.Printf("⚠️ Migration warning: %v", err)
	} else {
		utils.InfoLogger.Println("✅ Database migration completed.")
	}

	return db // ✅ Return the DB
}

// ✅ ConnectTestDB initializes a separate test DB
func ConnectTestDB() {
	utils.InitLogger()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=1016 dbname=shetiseva_test port=5432 sslmode=disable"
		utils.InfoLogger.Println("⚠️ TEST_DATABASE_URL not set. Using fallback test DSN.")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.ErrorLogger.Fatalf("❌ Failed to connect to test database: %v", err)
	}

	TestDB = db
	utils.InfoLogger.Println("✅ Connected to test database.")

	// Optional: Auto-migrate for test DB too
	err = db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.CartItem{},
		&models.Wishlist{},
		&models.Order{},
		&models.ContactMessage{},
	)
	if err != nil {
		utils.ErrorLogger.Printf("⚠️ Test DB migration warning: %v", err)
	} else {
		utils.InfoLogger.Println("✅ Test database migration completed.")
	}
}
