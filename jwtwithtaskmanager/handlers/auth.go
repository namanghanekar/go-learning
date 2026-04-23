package handlers

import (
	"strings"
	"time"

	"jwtwithtaskmanager/db"
	"jwtwithtaskmanager/models"
	"jwtwithtaskmanager/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// REGISTER → send OTP + store HASHED password
func RegisterHandler(c *gin.Context) {

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// 🔥 Trim input (IMPORTANT)
	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)

	// 🔥 Hash password BEFORE storing
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	// generate OTP
	otp := utils.GenerateOTP()

	// store OTP + HASH in Redis
	utils.RDB.Set(utils.Ctx, input.Username+":otp", otp, time.Minute*5)
	utils.RDB.Set(utils.Ctx, input.Username+":pass", string(hash), time.Minute*5)

	// send OTP
	utils.SendOTPEmail(input.Username, otp)

	c.JSON(200, gin.H{"message": "OTP sent to email"})
}

// VERIFY OTP → create user
func VerifyOTP(c *gin.Context) {

	var input struct {
		Username string `json:"username"`
		OTP      string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	input.Username = strings.TrimSpace(input.Username)

	// get OTP
	storedOTP, err := utils.RDB.Get(utils.Ctx, input.Username+":otp").Result()
	if err != nil {
		c.JSON(400, gin.H{"error": "OTP expired or not found"})
		return
	}

	if storedOTP != input.OTP {
		c.JSON(401, gin.H{"error": "Invalid OTP"})
		return
	}

	// get HASHED password
	hash, err := utils.RDB.Get(utils.Ctx, input.Username+":pass").Result()
	if err != nil {
		c.JSON(400, gin.H{"error": "Session expired, register again"})
		return
	}

	// create user
	user := models.User{
		Username:   input.Username,
		Password:   hash, // already hashed
		IsVerified: true,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	// cleanup Redis
	utils.RDB.Del(utils.Ctx, input.Username+":otp")
	utils.RDB.Del(utils.Ctx, input.Username+":pass")

	c.JSON(200, gin.H{"message": "User registered successfully"})
}

// LOGIN
func LoginHandler(c *gin.Context) {

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// 🔥 Trim input (CRITICAL)
	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)

	// find user
	if err := db.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// check verification
	if !user.IsVerified {
		c.JSON(401, gin.H{"error": "Verify OTP first"})
		return
	}

	// 🔥 Compare password correctly
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid password"})
		return
	}

	// generate JWT
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
