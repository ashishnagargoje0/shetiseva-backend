package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestCreateProductUnauthorized(t *testing.T) {
    productPayload := map[string]interface{}{
        "name":        "Test Product",
        "description": "Test Description",
        "price":       99.99,
        "category_id": 1,
    }
    body, _ := json.Marshal(productPayload)

    req, _ := http.NewRequest("POST", "/api/products", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    // Should fail because no auth token sent
    assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

// For authorized tests, you need to generate a JWT token from login and pass in Authorization header
// You can write helper functions to do that.
