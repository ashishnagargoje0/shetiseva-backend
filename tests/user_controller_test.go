package tests

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/ashishnagargoje0/shetiseva-backend/controllers"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
    router := gin.Default()
    router.GET("/users", controllers.GetUsers)

    req, _ := http.NewRequest("GET", "/users", nil)
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}
