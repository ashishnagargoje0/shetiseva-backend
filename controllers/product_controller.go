package controllers

import (
	"context"
	"net/http"

	"github.com/ashishnagargoje0/backend/database"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllProducts(c *gin.Context) {
	collection := database.GetCollection("products")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	collection := database.GetCollection("products")
	var product models.Product
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func GetAllCategories(c *gin.Context) {
	collection := database.GetCollection("categories")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer cursor.Close(context.TODO())

	var categories []models.Category
	if err := cursor.All(context.TODO(), &categories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetCategoryProducts(c *gin.Context) {
	slug := c.Param("slug")

	// Get category ID from slug
	catColl := database.GetCollection("categories")
	var cat models.Category
	err := catColl.FindOne(context.TODO(), bson.M{"slug": slug}).Decode(&cat)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Get products by category
	prodColl := database.GetCollection("products")
	cursor, err := prodColl.Find(context.TODO(), bson.M{"category_id": cat.Slug})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse product list"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProductFilters(c *gin.Context) {
	filters := gin.H{
		"brands": []string{"BASF", "Bayer", "Mahindra"},
		"categories": []string{"Seeds", "Fertilizer", "Tools"},
		"priceRange": []int{0, 5000},
	}
	c.JSON(http.StatusOK, filters)
}

