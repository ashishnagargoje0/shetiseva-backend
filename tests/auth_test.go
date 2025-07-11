package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/ashishnagargoje0/backend/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func cleanupTestUsers() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection := config.DB.Collection("users")

	// Remove test users by email and phone used in tests
	userCollection.DeleteMany(ctx, bson.M{
		"$or": []bson.M{
			{"email": "testuser1@example.com"},
			{"phone": "9999999999"},
			{"phone": "9876543210"},
			{"phone": "9876543211"},
		},
	})
}

func setupRouter() *gin.Engine {
	config.ConnectMongoDB()
	controllers.InitAuthCollection()
	cleanupTestUsers()

	r := gin.Default()
	r.POST("/auth/signup", controllers.Signup)
	r.POST("/auth/login", controllers.Login)
	r.POST("/auth/forgot-password", controllers.ForgotPassword)
	r.POST("/auth/register-phone", controllers.RegisterPhone)
	r.POST("/auth/verify-otp", controllers.VerifyOTP)
	return r
}

func TestSignupAndLogin(t *testing.T) {
	r := setupRouter()

	// 1. Signup
	t.Run("Signup", func(t *testing.T) {
		payload := map[string]string{
			"name":     "Test User",
			"email":    "testuser1@example.com",
			"password": "test1234",
			"phone":    "9876543210",
		}
		body, _ := json.Marshal(payload)
		req := createJSONRequest("POST", "/auth/signup", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Contains(t, []int{200, 201, 409}, resp.Code)
	})

	// 2. Duplicate Signup
	t.Run("Duplicate Signup", func(t *testing.T) {
		payload := map[string]string{
			"name":     "Test User",
			"email":    "testuser1@example.com",
			"password": "test1234",
			"phone":    "9876543210",
		}
		body, _ := json.Marshal(payload)
		req := createJSONRequest("POST", "/auth/signup", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 409, resp.Code)
	})

	// 3. Login with correct password
	t.Run("Login Success", func(t *testing.T) {
		payload := map[string]string{
			"email":    "testuser1@example.com",
			"password": "test1234",
		}
		body, _ := json.Marshal(payload)
		req := createJSONRequest("POST", "/auth/login", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})

	// 4. Login with wrong password
	t.Run("Login Wrong Password", func(t *testing.T) {
		payload := map[string]string{
			"email":    "testuser1@example.com",
			"password": "wrongpass",
		}
		body, _ := json.Marshal(payload)
		req := createJSONRequest("POST", "/auth/login", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 401, resp.Code)
	})

	// 5. Invalid email format
	t.Run("Invalid Email Format", func(t *testing.T) {
		payload := map[string]string{
			"email":    "invalid-email",
			"password": "test1234",
			"name":     "Invalid Email",
			"phone":    "9876543211",
		}
		body, _ := json.Marshal(payload)
		req := createJSONRequest("POST", "/auth/signup", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 400, resp.Code)
	})

	// 6. Forgot Password
	t.Run("Forgot Password", func(t *testing.T) {
		payload := map[string]string{
			"email": "testuser1@example.com",
		}
		body, _ := json.Marshal(payload)
		req := createJSONRequest("POST", "/auth/forgot-password", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})

	// 7. OTP Registration
	t.Run("Register Phone", func(t *testing.T) {
		payload := map[string]string{
			"phone": "9999999999",
			"name":  "OTP User",
		}
		body, _ := json.Marshal(payload)
		req := createJSONRequest("POST", "/auth/register-phone", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})

	// 8. OTP Verification
	t.Run("Verify OTP", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user struct {
			OTP string `bson:"otp"`
		}
		err := config.DB.Collection("users").FindOne(ctx, bson.M{"phone": "9999999999"}).Decode(&user)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.OTP)

		payload := map[string]string{
			"phone": "9999999999",
			"otp":   user.OTP,
		}
		body, _ := json.Marshal(payload)
		req := createJSONRequest("POST", "/auth/verify-otp", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	})
}

// Helper function to create JSON requests
func createJSONRequest(method, url string, body []byte) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
