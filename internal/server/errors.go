package server

import (
	"log/slog"
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
					slog.Error("App error occurred",
						slog.Int("code", appErr.Code),
						slog.String("message", appErr.Message),
						slog.Any("error", appErr.Err),
					)
				}
				c.JSON(appErr.Code, gin.H{"error": appErr.Message})
				return
			}

			slog.Error("Unexpected internal error occurred", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
		}
	}
}
