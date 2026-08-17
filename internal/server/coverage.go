package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) checkCoverageHandler(c *gin.Context) {
	postcode := c.Query("postcode")
	if postcode == "" {
		_ = c.Error(NewAppError(http.StatusBadRequest, "Query parameter 'postcode' is required", nil))
		return
	}

	covered, cleanPostcode, message := s.coverageService.CheckCoverage(postcode)
	c.JSON(http.StatusOK, gin.H{
		"postcode": cleanPostcode,
		"covered":  covered,
		"message":  message,
	})
}
