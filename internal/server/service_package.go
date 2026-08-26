package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (s *Server) listServicesHandler(c *gin.Context) {
	services, err := s.servicePackageService.ListServices(c.Request.Context())
	if err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to retrieve services", err))
		return
	}

	c.JSON(http.StatusOK, services)
}

func (s *Server) getServiceBySlugHandler(c *gin.Context) {
	slug := c.Param("slug")
	service, err := s.servicePackageService.GetServiceBySlug(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(NewAppError(http.StatusNotFound, "Service not found", err))
			return
		}
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to retrieve service", err))
		return
	}

	c.JSON(http.StatusOK, service)
}

