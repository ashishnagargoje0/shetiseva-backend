package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	AuthRoutes(router)
	UserRoutes(router)
	ProductRoutes(router)
	CartRoutes(router)
	OrderRoutes(router)
	CategoryRoutes(router)
	WishlistRoutes(router)
	ContactRoutes(router)
	AdminRoutes(router) // Make sure this function exists in contact_routes.go
}
