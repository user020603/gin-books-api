package utils

import (
	"github.com/gin-gonic/gin"
)

// JSONResponse sends a JSON response with the given status code and data.
func JSONResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// ErrorResponse sends an error response with the given status code and error message.
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}
