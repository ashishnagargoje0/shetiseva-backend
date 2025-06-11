package middleware

import (
	"net/http"

	"github.com/ashishnagargoje0/shetiseva-backend/config"
	"github.com/ashishnagargoje0/shetiseva-backend/models"
	"github.com/ashishnagargoje0/shetiseva-backend/utils"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			utils.HandleError(c, http.StatusUnauthorized, "Unauthorized: Missing user ID")
			c.Abort()
			return
		}

		// Ensure correct type
		userID, ok := userIDRaw.(uint)
		if !ok {
			utils.HandleError(c, http.StatusUnauthorized, "Unauthorized: Invalid user ID type")
			c.Abort()
			return
		}

		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			utils.HandleError(c, http.StatusUnauthorized, "Unauthorized: User not found")
			c.Abort()
			return
		}

		if user.Role != "admin" {
			utils.HandleError(c, http.StatusForbidden, "Access denied: Admins only")
			c.Abort()
			return
		}

		c.Next()
	}
}
