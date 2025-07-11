package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/ashishnagargoje0/backend/database"
	"github.com/ashishnagargoje0/backend/middlewares"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouterForAPI() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// ✅ Initialize MongoDB
	database.ConnectDB()

	r := gin.Default()

	// Public routes
	r.GET("/products", controllers.GetAllProducts)
	r.GET("/product/:id", controllers.GetProductByID)
	r.GET("/categories", controllers.GetAllCategories)
	r.GET("/category/:slug", controllers.GetCategoryProducts)
	r.GET("/marketplace/items", controllers.GetAllMarketplaceItems)
	r.GET("/marketplace/item/:id", controllers.GetMarketplaceItemByID)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.POST("/wishlist/add", controllers.AddToWishlist)
		auth.GET("/wishlist", controllers.GetWishlist)
		auth.DELETE("/wishlist/:id", controllers.RemoveFromWishlist)

		auth.POST("/compare/add", controllers.AddToCompare)
		auth.GET("/compare/view", controllers.GetCompareList)

		auth.POST("/cart/add", controllers.AddToCart)
		auth.GET("/cart/view/:user_id", controllers.ViewCart)
		auth.DELETE("/cart/remove", controllers.RemoveFromCart)
		auth.PUT("/cart/update-qty", controllers.UpdateCartQuantity)
		auth.POST("/checkout", controllers.Checkout)
	}

	// Auth routes
	r.POST("/auth/login", controllers.Login)

	return r
}

func getJWTToken(r *gin.Engine) (string, string) {
	loginBody := map[string]string{
		"email":    "admin@shetiseva.in",
		"password": "admin123",
	}
	jsonBody, _ := json.Marshal(loginBody)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		panic("❌ Failed to login and get JWT token for tests")
	}

	var respBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		panic("❌ Failed to decode login response")
	}

	token := respBody["token"].(string)
	userID := respBody["user"].(map[string]interface{})["id"].(string)

	return "Bearer " + token, userID
}

func TestProductAndMarketplaceRoutes(t *testing.T) {
	r := setupRouterForAPI()
	token, userID := getJWTToken(r)

	// ✅ Insert test product with real ObjectID
	testProduct := models.Product{
		ID:    primitive.NewObjectID(),
		Name:  "Test Product",
		Price: 200,
	}
	_, err := database.ProductCollection.InsertOne(context.Background(), testProduct)
	if err != nil {
		t.Fatalf("❌ Failed to insert test product: %v", err)
	}
	testProductID := testProduct.ID.Hex()

	// ✅ Clean up after all tests
	defer func() {
		database.ProductCollection.DeleteOne(context.Background(), bson.M{"_id": testProduct.ID})
	}()

	t.Run("Get All Products", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/products", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("Get Categories", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/categories", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("Get Marketplace Items", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/marketplace/items", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("Add to Wishlist", func(t *testing.T) {
		body := map[string]string{"product_id": testProductID}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/wishlist/add", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Contains(t, []int{200, 409}, resp.Code)
	})

	t.Run("View Wishlist", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wishlist", nil)
		req.Header.Set("Authorization", token)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("Add to Compare", func(t *testing.T) {
		body := map[string]string{"product_id": testProductID}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/compare/add", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Contains(t, []int{200, 409}, resp.Code)
	})

	t.Run("View Compare List", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/compare/view", nil)
		req.Header.Set("Authorization", token)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("Add to Cart", func(t *testing.T) {
		body := map[string]interface{}{
			"product_id": testProductID,
			"quantity":   1,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/cart/add", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("View Cart", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cart/view/"+userID, nil)
		req.Header.Set("Authorization", token)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("Remove from Cart", func(t *testing.T) {
		body := map[string]string{
			"user_id":    userID,
			"product_id": testProductID,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest("DELETE", "/cart/remove", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Contains(t, []int{200, 404}, resp.Code)
	})

	t.Run("Update Cart Quantity", func(t *testing.T) {
		body := map[string]interface{}{
			"user_id":    userID,
			"product_id": testProductID,
			"quantity":   2,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/cart/update-qty", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Contains(t, []int{200, 404}, resp.Code)
	})
}
