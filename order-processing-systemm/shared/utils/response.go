package utils

import "github.com/gin-gonic/gin"

func SuccessResponse(
	c *gin.Context,
	message string,
	data interface{},
) {

	c.JSON(200, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(
	c *gin.Context,
	message string,
) {

	c.JSON(400, gin.H{
		"success": false,
		"message": message,
	})
}
