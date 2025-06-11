package middleware

import (
	"net/http"
	"strings"

	"github.com/ashishnagargoje0/shetiseva-backend/utils"
	"github.com/gin-gonic/gin"
)

// Auth middleware validates JWT and adds user info to context
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.HandleError(c, http.StatusUnauthorized, "Unauthorized: Missing or invalid Authorization header")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// ✅ FIXED: Correct function name
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			utils.HandleError(c, http.StatusUnauthorized, "Unauthorized: Invalid or expired token")
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok || userIDFloat <= 0 {
			utils.HandleError(c, http.StatusUnauthorized, "Unauthorized: Invalid user ID in token")
			c.Abort()
			return
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			utils.HandleError(c, http.StatusUnauthorized, "Unauthorized: Invalid email in token")
			c.Abort()
			return
		}

		c.Set("user_id", uint(userIDFloat))
		c.Set("email", email)

		c.Next()
	}
}
