package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks = []Task{
	{ID: 1, Title: "Go learning", Done: false},
	{ID: 2, Title: "Build API", Done: false},
}

var nextID = 3

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()

		log.Printf("Response Status: %d", c.Writer.Status())
	}
}
func getTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}
func getTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	for _, t := range tasks {
		if t.ID == id {
			c.JSON(200, t)
			return
		}
	}

	c.JSON(404, gin.H{"error": "Task not found"})
}
func createTask(c *gin.Context) {
	var newTask Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	newTask.ID = nextID
	nextID++

	tasks = append(tasks, newTask)

	c.JSON(201, newTask)
}
func updateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	var updated Task
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Title = updated.Title
			tasks[i].Done = updated.Done
			c.JSON(200, tasks[i])
			return
		}
	}

	c.JSON(404, gin.H{"error": "Task not found"})
}
func deleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(200, gin.H{"message": "Task deleted"})
			return
		}
	}

	c.JSON(404, gin.H{"error": "Task not found"})
}

func main() {
	r := gin.Default()

	// Apply middleware
	r.Use(Logger())

	// Routes
	r.GET("/tasks", getTasks)
	r.GET("/tasks/:id", getTask)
	r.POST("/tasks", createTask)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)

	r.Run(":8080")
}
