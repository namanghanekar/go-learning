package handlers

import (
	"jwtwithtaskmanager/db"
	"jwtwithtaskmanager/models"
	"jwtwithtaskmanager/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "User exists"})
		return
	}

	c.JSON(201, gin.H{"message": "User registered"})
}

func LoginHandler(c *gin.Context) {
	var input models.User
	var dbUser models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.DB.Where("username = ?", input.Username).First(&dbUser).Error; err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(input.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Wrong password"})
		return
	}

	token, _ := utils.GenerateJWT(dbUser.ID)

	c.JSON(200, gin.H{"token": token})
}
