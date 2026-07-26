package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id,omitempty"`
	Details   interface{} `json:"details,omitempty"`
}

func Error(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error": APIError{
			Code:      code,
			Message:   message,
			RequestID: c.GetString("request_id"),
		},
	})
}

func ValidationError(c *gin.Context, details interface{}) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error": APIError{
			Code:      "validation_error",
			Message:   "La solicitud contiene datos inválidos",
			RequestID: c.GetString("request_id"),
			Details:   details,
		},
	})
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}
