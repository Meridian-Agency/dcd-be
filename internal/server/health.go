package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) healthHandler(c *gin.Context) {
	health := s.db.Health()
	if health["status"] != "up" {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Database connection is unhealthy", nil))
		return
	}
	c.JSON(http.StatusOK, health)
}

func (s *Server) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
