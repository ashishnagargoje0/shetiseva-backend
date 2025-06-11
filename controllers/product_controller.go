package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"github.com/ashishnagargoje0/shetiseva-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var productDB *gorm.DB

func InitProductController(db *gorm.DB) {
	productDB = db
}

// ------------------- CREATE PRODUCT -------------------
func CreateProduct(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	description := strings.TrimSpace(c.PostForm("description"))
	priceStr := c.PostForm("price")
	quantityStr := c.PostForm("quantity")
	categoryIDStr := c.PostForm("category_id")

	// Validate basic required fields
	if name == "" || description == "" || priceStr == "" || quantityStr == "" || categoryIDStr == "" {
		utils.HandleError(c, http.StatusBadRequest, "All fields (name, description, price, quantity, category_id) are required")
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		utils.HandleError(c, http.StatusBadRequest, "Price must be a positive number")
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 0 {
		utils.HandleError(c, http.StatusBadRequest, "Quantity must be a non-negative integer")
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil || categoryID <= 0 {
		utils.HandleError(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	// Image upload
	file, err := c.FormFile("image")
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Image file is required")
		return
	}

	safeFileName := filepath.Base(file.Filename)
	imagePath := filepath.Join("uploads", safeFileName)
	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to save image")
		return
	}

	// Create product
	product := models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		CategoryID:  uint(categoryID),
		ImageURL:    "/" + imagePath,
	}

	if err := productDB.Create(&product).Error; err != nil {
		fmt.Println("❌ DB Error while creating product:", err)
		utils.HandleError(c, http.StatusInternalServerError, "Database error while creating product")
		return
	}

	c.JSON(http.StatusCreated, product)
}

// ------------------- GET ALL PRODUCTS -------------------
func GetProducts(c *gin.Context) {
	var products []models.Product

	// Query params
	search := c.Query("search")
	categoryID := c.Query("category_id")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")
	sort := c.DefaultQuery("sort", "created_at_desc")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	query := productDB.Preload("Category")

	if search != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			query = query.Where("price >= ?", minPrice)
		}
	}
	if maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			query = query.Where("price <= ?", maxPrice)
		}
	}

	// Sorting
	switch sort {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "name_asc":
		query = query.Order("name ASC")
	case "name_desc":
		query = query.Order("name DESC")
	default:
		query = query.Order("created_at DESC")
	}

	var total int64
	query.Model(&models.Product{}).Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       products,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

// ------------------- GET SINGLE PRODUCT -------------------
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := productDB.Preload("Category").First(&product, id).Error; err != nil {
		utils.HandleError(c, http.StatusNotFound, "Product not found")
		return
	}
	c.JSON(http.StatusOK, product)
}

// ------------------- UPDATE PRODUCT -------------------
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := productDB.First(&product, id).Error; err != nil {
		utils.HandleError(c, http.StatusNotFound, "Product not found")
		return
	}

	if name := strings.TrimSpace(c.PostForm("name")); name != "" {
		product.Name = name
	}
	if description := strings.TrimSpace(c.PostForm("description")); description != "" {
		product.Description = description
	}
	if priceStr := c.PostForm("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil && price > 0 {
			product.Price = price
		}
	}
	if quantityStr := c.PostForm("quantity"); quantityStr != "" {
		if quantity, err := strconv.Atoi(quantityStr); err == nil && quantity >= 0 {
			product.Quantity = quantity
		}
	}
	if categoryIDStr := c.PostForm("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.Atoi(categoryIDStr); err == nil && categoryID > 0 {
			product.CategoryID = uint(categoryID)
		}
	}
	if file, err := c.FormFile("image"); err == nil {
		safeFileName := filepath.Base(file.Filename)
		imagePath := filepath.Join("uploads", safeFileName)
		if err := c.SaveUploadedFile(file, imagePath); err == nil {
			product.ImageURL = "/" + imagePath
		}
	}

	if err := productDB.Save(&product).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to update product")
		return
	}

	c.JSON(http.StatusOK, product)
}

// ------------------- DELETE PRODUCT -------------------
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := productDB.Delete(&models.Product{}, id).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to delete product")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// ------------------- GET LOW STOCK PRODUCTS -------------------
func GetLowStockProducts(c *gin.Context) {
	var lowStockProducts []models.Product
	if err := productDB.Where("quantity < ?", 10).Preload("Category").Find(&lowStockProducts).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to fetch low stock products")
		return
	}
	c.JSON(http.StatusOK, lowStockProducts)
}
