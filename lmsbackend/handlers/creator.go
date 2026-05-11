package handlers

import (
	"strings"

	"lmsbackend/config"
	"lmsbackend/constants"
	"lmsbackend/models"
	"lmsbackend/utils"

	"github.com/gin-gonic/gin"
)

func RegisterCreator(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Platform    string `json:"platform"`
		Followers   string `json:"followers"`
		ProfileLink string `json:"profileLink"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	input.Email = strings.TrimSpace(input.Email)
	input.Password = strings.TrimSpace(input.Password)
	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Password hashing failed"})
		return
	}
	user := models.User{
		Email:      input.Email,
		Password:   hash,
		Role:       constants.RoleCreator,
		IsVerified: true,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "User creation failed"})
		return
	}
	profile := models.CreatorProfile{
		UserID:      user.ID,
		Name:        input.Name,
		Platform:    input.Platform,
		Followers:   input.Followers,
		ProfileLink: input.ProfileLink,
	}
	config.DB.Create(&profile)
	c.JSON(200, gin.H{"message": "Creator registered successfully"})
}
