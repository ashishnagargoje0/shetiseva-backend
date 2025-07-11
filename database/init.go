package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitDatabase sets up indexes and initial database state
func InitDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := config.DB

	// Ensure unique index on "email" in users collection
	userCol := db.Collection("users")
	emailIndex := mongoIndex("email", true)
	if _, err := userCol.Indexes().CreateOne(ctx, emailIndex); err != nil {
		log.Fatalf("❌ Failed to create user email index: %v", err)
	}
	fmt.Println("✅ Created unique index on users.email")

	// Add more indexes for performance (example: phone, contact timestamp, etc.)
	contactCol := db.Collection("contacts")
	timestampIndex := mongoIndex("created_at", false)
	if _, err := contactCol.Indexes().CreateOne(ctx, timestampIndex); err != nil {
		log.Printf("⚠️ Failed to create index on contacts.created_at: %v", err)
	} else {
		fmt.Println("✅ Index on contacts.created_at")
	}

	// Example: add unique index on phone for KYC, if needed
	kycCol := db.Collection("kyc_docs")
	phoneIndex := mongoIndex("phone", true)
	if _, err := kycCol.Indexes().CreateOne(ctx, phoneIndex); err != nil {
		log.Printf("⚠️ KYC phone index not created: %v", err)
	} else {
		fmt.Println("✅ KYC phone unique index created")
	}
}

// mongoIndex is a helper to define a MongoDB index
func mongoIndex(field string, unique bool) mongo.IndexModel {
	return mongo.IndexModel{
		Keys: bson.D{{Key: field, Value: 1}},
		Options: &options.IndexOptions{
			Unique: &unique,
		},
	}
}

// GetCollection returns a MongoDB collection reference from the configured DB
func GetCollection(name string) *mongo.Collection {
	return config.DB.Collection(name)
}
