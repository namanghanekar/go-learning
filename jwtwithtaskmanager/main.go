package main

import (
	"jwtwithtaskmanager/db"
	"jwtwithtaskmanager/handlers"
	"jwtwithtaskmanager/middleware"
	"jwtwithtaskmanager/models"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	db.DB.AutoMigrate(&models.User{}, &models.Task{})
	r := gin.Default()
	r.POST("/register", handlers.RegisterHandler)
	r.POST("/login", handlers.LoginHandler)
	protected := r.Group("/")
	protected.Use(middleware.JWTMiddleware())
	{
		protected.POST("/task", handlers.CreateTask)
		protected.GET("/tasks", handlers.GetTasks)
		protected.PUT("/task/:id", handlers.UpdateTask)
		protected.DELETE("/task/:id", handlers.DeleteTask)
	}
	r.Run(":8080")
}
