package database

import (
	"log"

	"github.com/ashishnagargoje0/backend/config"
	"go.mongodb.org/mongo-driver/mongo"
)

// Global collection variables
var (
	UserCollection     *mongo.Collection
	ProductCollection  *mongo.Collection
	CartCollection     *mongo.Collection
	WishlistCollection *mongo.Collection
	CompareCollection  *mongo.Collection
	ContactCollection  *mongo.Collection
	KYCCollection      *mongo.Collection
	OrderCollection    *mongo.Collection
)

// ConnectDB assigns MongoDB collections after config.DB is connected
func ConnectDB() {
	if config.DB == nil {
		log.Fatal("❌ MongoDB not initialized. Call config.ConnectMongoDB() first.")
	}

	// Assign collections from config.DB
	UserCollection = config.DB.Collection("users")
	ProductCollection = config.DB.Collection("products")
	CartCollection = config.DB.Collection("cart")
	WishlistCollection = config.DB.Collection("wishlist")
	CompareCollection = config.DB.Collection("compare")
	ContactCollection = config.DB.Collection("contacts")
	KYCCollection = config.DB.Collection("kyc_docs")
	OrderCollection = config.DB.Collection("orders")

	log.Println("✅ MongoDB collections assigned.")
}
