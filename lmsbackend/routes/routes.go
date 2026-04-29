package routes

import (
	"lmsbackend/handlers"
	"lmsbackend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Public routes
	r.POST("/register/creator", handlers.RegisterCreator)
	r.POST("/register/brand", handlers.RegisterBrand)
	r.POST("/login", handlers.Login)

	// Protected routes (JWT required)
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	api.GET("/profile", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")

		c.JSON(200, gin.H{
			"user_id": userID,
			"role":    role,
		})
	})
}
