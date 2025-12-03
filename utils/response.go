package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(200, BaseResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Error:   nil,
	})
}

func Error(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, BaseResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
		Error:   err,
	})
}

func SuccessLang(c *gin.Context, message string, idData interface{}, enData interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data": gin.H{
			"id": idData,
			"en": enData,
		},
	})
}
