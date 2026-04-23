package main

import (
	"jwtwithtaskmanager/db"
	"jwtwithtaskmanager/handlers"
	"jwtwithtaskmanager/middleware"
	"jwtwithtaskmanager/models"
	"jwtwithtaskmanager/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()
	utils.InitRedis()
	db.DB.AutoMigrate(&models.User{}, &models.Task{})
	r := gin.Default()
	// Public
	r.POST("/register", handlers.RegisterHandler)
	r.POST("/verify-otp", handlers.VerifyOTP)
	r.POST("/login", handlers.LoginHandler)

	// Protected
	protected := r.Group("/")
	protected.Use(middleware.JWTMiddleware())
	{
		protected.POST("/task", handlers.CreateTask)
		protected.GET("/tasks", handlers.GetTasks)
	}

	r.Run(":8080")
}
