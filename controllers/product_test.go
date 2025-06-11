package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"github.com/ashishnagargoje0/shetiseva-backend/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup test DB and controller
	db := testutils.SetupTestDB()
	InitProductController(db)

	// Optional: Seed one product into test DB to ensure response
	db.Create(&models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Quantity:    10,
		CategoryID:  1,
		ImageURL:    "/uploads/test.jpg",
	})

	// Setup router
	router := gin.Default()
	router.GET("/products", GetProducts)

	// Perform request
	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert status and body content
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Product")
}
