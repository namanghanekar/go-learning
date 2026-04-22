package handlers

import (
	"jwtwithtaskmanager/db"
	"jwtwithtaskmanager/models"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	task.UserID = uint(userID.(float64))
	if err := db.DB.Create(&task).Error; err != nil {
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	c.JSON(200, gin.H{"message": "Task created"})
}
func GetTasks(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var tasks []models.Task
	db.DB.Where("user_id = ?", userID).Find(&tasks)
	c.JSON(200, tasks)
}
func UpdateTask(c *gin.Context) {

	// get user_id from middleware
	userID, _ := c.Get("user_id")

	// get task id from URL
	id := c.Param("id")

	var task models.Task

	// bind JSON
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// update only if task belongs to user
	err := db.DB.Model(&models.Task{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("title", task.Title).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Task updated"})
}
func DeleteTask(c *gin.Context) {

	userID, _ := c.Get("user_id")

	id := c.Param("id")

	err := db.DB.Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Task{}).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Delete failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Task deleted"})
}
