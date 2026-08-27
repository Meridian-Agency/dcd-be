package server

import (
	"net/http"

	"dcd-be/internal/service"

	"github.com/gin-gonic/gin"
)

type CreateReservationDTO struct {
	Name    string `json:"name" binding:"required"`
	Contact string `json:"contact" binding:"required"`
	Service string `json:"service" binding:"required"`
	Date    string `json:"date" binding:"required,datetime=2006-01-02T15:04:05Z07:00"` // Validates RFC3339 layout
	Message string `json:"message"`
}

func (s *Server) createReservationHandler(c *gin.Context) {
	var input CreateReservationDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, err.Error(), err))
		return
	}

	reservationInput := service.CreateReservationInput{
		Name:    input.Name,
		Contact: input.Contact,
		Service: input.Service,
		Date:    input.Date,
		Message: input.Message,
	}

	res, err := s.reservationService.CreateReservation(c.Request.Context(), reservationInput)
	if err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to store reservation", err))
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (s *Server) listReservationsHandler(c *gin.Context) {
	publicList, err := s.reservationService.ListReservationsPublic(c.Request.Context())
	if err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to retrieve reservations", err))
		return
	}

	c.JSON(http.StatusOK, publicList)
}

func (s *Server) listAdminReservationsHandler(c *gin.Context) {
	adminList, err := s.reservationService.ListReservationsAdmin(c.Request.Context())
	if err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to retrieve reservations", err))
		return
	}

	c.JSON(http.StatusOK, adminList)
}
