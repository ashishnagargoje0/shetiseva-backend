package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/ashishnagargoje0/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// SeedDatabase seeds initial data (admin user, districts, etc.)
func SeedDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := config.DB

	// ✅ 1. Seed Admin User
	adminCol := db.Collection("users")
	adminEmail := "admin@shetiseva.in"
	count, _ := adminCol.CountDocuments(ctx, bson.M{"email": adminEmail})
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.User{
			Name:     "Super Admin",
			Email:    adminEmail,
			Password: string(hashedPassword),
			Role:     "admin",
		}
		if _, err := adminCol.InsertOne(ctx, admin); err != nil {
			log.Printf("⚠️ Failed to insert admin user: %v", err)
		} else {
			fmt.Println("✅ Default admin user created: admin@shetiseva.in / admin123")
		}
	} else {
		fmt.Println("ℹ️ Admin user already exists.")
	}

	// ✅ 2. Seed District List (for dropdown)
	districtCol := db.Collection("districts")
	districts := []string{"Pune", "Nagpur", "Nashik", "Solapur", "Aurangabad"}

	for _, d := range districts {
		filter := bson.M{"name": d}
		update := bson.M{"$setOnInsert": bson.M{"name": d, "created_at": time.Now()}}
		opts := options.Update().SetUpsert(true)

		if _, err := districtCol.UpdateOne(ctx, filter, update, opts); err != nil {
			log.Printf("⚠️ Failed to insert district: %s", d)
		}
	}
	fmt.Println("✅ Districts seeded.")

	// ✅ 3. Optional: Seed Configs or Initial Products (skip if already done)
}
