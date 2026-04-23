package handlers

import (
	"jwtwithtaskmanager/db"
	"jwtwithtaskmanager/models"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	var task models.Task
	c.ShouldBindJSON(&task)

	task.UserID = userID

	db.DB.Create(&task)

	c.JSON(200, gin.H{"message": "Task created"})
}

func GetTasks(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	var tasks []models.Task
	db.DB.Where("user_id = ?", userID).Find(&tasks)

	c.JSON(200, tasks)
}
