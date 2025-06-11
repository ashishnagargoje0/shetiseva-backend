package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ashishnagargoje0/shetiseva-backend/config"
	"github.com/ashishnagargoje0/shetiseva-backend/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	// Connect to the test database
	config.ConnectTestDB()
	testDB := config.TestDB

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize router and routes
	router = gin.Default()
	auth := router.Group("/api/auth")
	{
		auth.POST("/signup", controllers.Register) // ✅ FIXED: use Register instead of SignUp
		auth.POST("/login", controllers.Login)
	}

	// Run all tests
	code := m.Run()

	// Cleanup test user
	if testDB != nil {
		testDB.Exec("DELETE FROM users WHERE email = ?", "testuser@example.com")
	}

	// Exit with test code
	os.Exit(code)
}

func TestSignup(t *testing.T) {
	payload := map[string]string{
		"name":     "Test User",
		"email":    "testuser@example.com",
		"password": "password123",
	}

	body, err := json.Marshal(payload)
	assert.NoError(t, err, "Failed to marshal signup payload")

	req, err := http.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(body))
	assert.NoError(t, err, "Failed to create signup request")
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code, "Expected status 201 Created, got %d", resp.Code)
}

func TestLogin(t *testing.T) {
	payload := map[string]string{
		"email":    "testuser@example.com",
		"password": "password123",
	}

	body, err := json.Marshal(payload)
	assert.NoError(t, err, "Failed to marshal login payload")

	req, err := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	assert.NoError(t, err, "Failed to create login request")
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "Expected status 200 OK")

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read login response body")

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	assert.NoError(t, err, "Failed to unmarshal login response")

	assert.Contains(t, result, "token", "Response should contain 'token' field")
}
