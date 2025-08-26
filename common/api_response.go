package common

import "github.com/gin-gonic/gin"

func JSON(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, ApiResponse{
		Message: message,
		Data:    data,
	})
}