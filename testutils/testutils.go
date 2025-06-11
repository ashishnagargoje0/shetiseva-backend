package testutils

import (
	"log"

	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{}, &models.Wishlist{})
	if err != nil {
		log.Fatalf("❌ Auto migration failed: %v", err)
	}

	TestDB = db
	return db
}

func ResetTestDB(db *gorm.DB) {
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM categories")
	db.Exec("DELETE FROM wishlists")
}
