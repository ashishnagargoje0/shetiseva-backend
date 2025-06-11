package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashishnagargoje0/shetiseva-backend/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategory(t *testing.T) {
	// Use test DB
	ConnectTestDB()

	// Setup router (without middleware for isolated test)
	router := gin.Default()
	router.POST("/categories", controllers.CreateCategory)

	// Prepare request body
	category := map[string]string{"name": "Seeds"}
	body, _ := json.Marshal(category)

	req, err := http.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response status
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response
	var res map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)

	// Ensure structure has "data" and inside that, "name"
	data, ok := res["data"].(map[string]interface{})
	assert.True(t, ok, "response should contain 'data' object")

	assert.Equal(t, "Seeds", data["name"])

	// Clean up the inserted test data
	TestDB.Exec("DELETE FROM categories WHERE name = ?", "Seeds")
}
