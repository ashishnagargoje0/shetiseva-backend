package database


import (
	"context"
	"log"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RunMigrations applies manual migrations to MongoDB
func RunMigrations() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := config.DB

	// 🚀 Migration 1: Add default roles if not present
	rolesCol := db.Collection("roles")
	defaultRoles := []string{"user", "admin", "agent"}

	for _, role := range defaultRoles {
		filter := bson.M{"name": role}
		update := bson.M{"$setOnInsert": bson.M{"name": role, "created_at": time.Now()}}
		opts := options.Update().SetUpsert(true)

		_, err := rolesCol.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("⚠️ Failed to upsert role '%s': %v", role, err)
		} else {
			log.Printf("✅ Ensured role '%s' exists", role)
		}
	}

	// 🚀 Migration 2: Update all users missing the `role` field to "user"
	userCol := db.Collection("users")
	updateRole := bson.M{
		"$set": bson.M{"role": "user"},
	}
	res, err := userCol.UpdateMany(ctx, bson.M{"role": bson.M{"$exists": false}}, updateRole)
	if err != nil {
		log.Printf("⚠️ Failed to update missing roles: %v", err)
	} else {
		log.Printf("✅ Updated %d users with missing role", res.ModifiedCount)
	}

	// 🚀 Migration 3: Add `status: active` to existing orders if missing
	orderCol := db.Collection("orders")
	res2, err := orderCol.UpdateMany(ctx, bson.M{"status": bson.M{"$exists": false}}, bson.M{
		"$set": bson.M{"status": "pending"},
	})
	if err != nil {
		log.Printf("⚠️ Failed to update order statuses: %v", err)
	} else {
		log.Printf("✅ Updated %d orders with default status", res2.ModifiedCount)
	}
}
