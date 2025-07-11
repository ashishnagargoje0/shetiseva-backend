package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cartCollection *mongo.Collection

// Initialize cart collection
func InitCartCollection() {
	if config.DB == nil {
		panic("❌ MongoDB not initialized: config.DB is nil in InitCartCollection")
	}
	cartCollection = config.DB.Collection("cart")
}

// Internal helper: Get logged-in user's ObjectID from context
func getUserObjectID(c *gin.Context) (primitive.ObjectID, bool) {
	val, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return primitive.NilObjectID, false
	}

	email, ok := val.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token context"})
		return primitive.NilObjectID, false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return primitive.NilObjectID, false
	}

	return user.ID, true
}

// ✅ Add item to cart or increment quantity if already present
func AddToCart(c *gin.Context) {
	var item models.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, ok := getUserObjectID(c)
	if !ok {
		return
	}

	item.UserID = userID
	item.ID = primitive.NewObjectID()
	item.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": item.UserID, "product_id": item.ProductID}
	update := bson.M{
		"$inc": bson.M{"quantity": item.Quantity},
		"$setOnInsert": bson.M{
			"_id":        item.ID,
			"created_at": item.CreatedAt,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := cartCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to cart"})
}

// ✅ View all items in user's cart
func ViewCart(c *gin.Context) {
	userID, ok := getUserObjectID(c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	cursor, err := cartCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer cursor.Close(ctx)

	var items []models.CartItem
	if err := cursor.All(ctx, &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading cart items"})
		return
	}

	if len(items) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Cart is empty",
			"cart":    []models.CartItem{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart fetched successfully",
		"cart":    items,
	})
}


// ✅ Remove item by product_id
func RemoveFromCart(c *gin.Context) {
	var input struct {
		ProductID string `json:"product_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, ok := getUserObjectID(c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "product_id": input.ProductID}
	_, err := cartCollection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

// ✅ Update cart item quantity
func UpdateCartQuantity(c *gin.Context) {
	var input struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Quantity < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, ok := getUserObjectID(c)
	if !ok {
		return
	}

	// Convert product_id string to ObjectID
	productObjID, err := primitive.ObjectIDFromHex(input.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "product_id": productObjID}
	update := bson.M{"$set": bson.M{"quantity": input.Quantity}}

	result, err := cartCollection.UpdateOne(ctx, filter, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found or update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart quantity updated"})
}


// ✅ Checkout: creates order & clears cart
func Checkout(c *gin.Context) {
	userID, ok := getUserObjectID(c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	cursor, err := cartCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}
	defer cursor.Close(ctx)

	var cartItems []models.CartItem
	if err := cursor.All(ctx, &cartItems); err != nil || len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty or unreadable"})
		return
	}

	order := bson.M{
		"user_id":    userID,
		"items":      cartItems,
		"status":     "pending",
		"created_at": time.Now(),
	}

	orderCollection := config.DB.Collection("orders")
	_, err = orderCollection.InsertOne(ctx, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place order"})
		return
	}

	_, err = cartCollection.DeleteMany(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order placed but failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": order})
}

// ✅ View cart using query param for admin or external cases
func ViewCartByQueryParam(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	cursor, err := cartCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer cursor.Close(ctx)

	var items []models.CartItem
	if err := cursor.All(ctx, &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading data"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// ✅ Remove cart item by MongoDB _id
func RemoveCartItemByID(c *gin.Context) {
	itemID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := cartCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil || res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found or delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed successfully"})
}
