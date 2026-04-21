package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Model
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// In-memory data (temporary DB)
var users = []User{
	{ID: 1, Name: "Sreyansh", Age: 22},
	{ID: 2, Name: "Rahul", Age: 25},
}

var nextID = 3

// 🔹 GET all users
func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

// 🔹 GET user by ID
func getUser(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	for _, u := range users {
		if u.ID == id {
			c.JSON(http.StatusOK, u)
			return
		}
	}
}
func getUserByName(c *gin.Context) {

	name := c.Param("name")
	for _, u := range users {
		if u.Name == name {
			c.JSON(http.StatusOK, u)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

// 🔹 CREATE user
func createUser(c *gin.Context) {
	var newUser User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	newUser.ID = nextID
	nextID++

	users = append(users, newUser)

	c.JSON(http.StatusCreated, newUser)
}

// 🔹 UPDATE user
func updateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	for i, u := range users {
		if u.ID == id {
			users[i].Name = updatedUser.Name
			users[i].Age = updatedUser.Age

			c.JSON(http.StatusOK, users[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

// 🔹 DELETE user
func deleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func main() {
	r := gin.Default()

	// Routes
	r.GET("/users", getUsers)
	r.GET("/users/:id", getUser)
	r.GET("/users/name/:name", getUserByName)
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	r.Run(":8080")
}
