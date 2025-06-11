package tests

import (
	"log"
	"os"

	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func ConnectTestDB() {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=1016 dbname=shetiseva_test port=5432 sslmode=disable"
		log.Println("⚠️ TEST_DATABASE_URL not set, using fallback DSN")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to test DB: %v", err)
	}

	// Ping DB to verify connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get generic DB object: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Ping failed for test DB: %v", err)
	}

	// Migrate models for testing
	if err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Category{},
	); err != nil {
		log.Fatalf("❌ Failed to auto-migrate test DB: %v", err)
	}

	TestDB = db
	log.Println("✅ Connected and migrated test DB successfully.")
}
