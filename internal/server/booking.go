package server

import (
	"net/http"

	"dcd-be/internal/service"

	"github.com/gin-gonic/gin"
)

type CreateBookingDTO struct {
	CustomerName  string `json:"customer_name" binding:"required,min=2"`
	CustomerEmail string `json:"customer_email" binding:"required,email"`
	CustomerPhone string `json:"customer_phone" binding:"required"`
	LocationType  string `json:"location_type" binding:"required,oneof=STUDIO MOBILE"`
	Postcode      string `json:"postcode" binding:"required_if=LocationType MOBILE"`
	VehicleMake   string `json:"vehicle_make" binding:"required"`
	VehicleModel  string `json:"vehicle_model" binding:"required"`
	VehicleSize   string `json:"vehicle_size" binding:"required"`
	PreferredDate string `json:"preferred_date" binding:"required"`
	Notes         string `json:"notes"`
}

func (s *Server) createBookingHandler(c *gin.Context) {
	var input CreateBookingDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, err.Error(), err))
		return
	}

	bookingInput := service.CreateBookingInput{
		CustomerName:  input.CustomerName,
		CustomerEmail: input.CustomerEmail,
		CustomerPhone: input.CustomerPhone,
		LocationType:  input.LocationType,
		Postcode:      input.Postcode,
		VehicleMake:   input.VehicleMake,
		VehicleModel:  input.VehicleModel,
		VehicleSize:   input.VehicleSize,
		PreferredDate: input.PreferredDate,
		Notes:         input.Notes,
	}

	booking, err := s.bookingService.CreateBooking(c.Request.Context(), bookingInput)
	if err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to store booking", err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":        "Booking request received. DCD Studio will contact you within 30 minutes.",
		"reference_code": booking.ReferenceCode,
		"status":         booking.Status,
	})
}
