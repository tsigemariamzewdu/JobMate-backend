package utils

import "github.com/gin-gonic/gin"




func SuccessPayload(message string, data any) gin.H {
	return gin.H{
		"success": true,
		"message": message,
		"data":    data,
	}
}

func ErrorPayload(message string, details any) gin.H {
	return gin.H{
		"success": false,
		"message": message,
		"details": details,
	}
}