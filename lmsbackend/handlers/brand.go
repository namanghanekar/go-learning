package handlers

import (
	"strings"

	"lmsbackend/config"
	"lmsbackend/constants"
	"lmsbackend/models"
	"lmsbackend/utils"

	"github.com/gin-gonic/gin"
)

func RegisterBrand(c *gin.Context) {

	var input struct {
		CompanyName string `json:"companyName"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Website     string `json:"website"`
		GST         string `json:"gst"`
		Description string `json:"description"`
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
		Role:       constants.RoleBrand,
		IsVerified: true,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "User creation failed"})
		return
	}

	profile := models.BrandProfile{
		UserID:      user.ID,
		CompanyName: input.CompanyName,
		Website:     input.Website,
		GST:         input.GST,
		Description: input.Description,
	}

	config.DB.Create(&profile)

	c.JSON(200, gin.H{"message": "Brand registered successfully"})
}
