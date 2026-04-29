package main

import (
	"lmsbackend/config"
	"lmsbackend/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	// Load configs
	config.LoadConfig()

	// Connect DB + Redis
	config.ConnectDB()
	config.ConnectRedis()

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	routes.SetupRoutes(r)

	r.Run(":8080")
}
