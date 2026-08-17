package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := err.(*AppError); ok {
				if appErr.Err != nil {
					log.Printf("[ERROR] %d: %s | %v", appErr.Code, appErr.Message, appErr.Err)
				}
				c.JSON(appErr.Code, gin.H{"error": appErr.Message})
				return
			}

			log.Printf("[UNEXPECTED ERROR] %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		}
	}
}
