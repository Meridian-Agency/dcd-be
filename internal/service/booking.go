package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"dcd-be/internal/database"
)

type BookingService interface {
	CreateBooking(ctx context.Context, input CreateBookingInput) (*database.Booking, error)
}

type bookingService struct {
	db database.Service
}

func NewBookingService(db database.Service) BookingService {
	return &bookingService{db: db}
}

type CreateBookingInput struct {
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	LocationType  string
	Postcode      string
	VehicleMake   string
	VehicleModel  string
	VehicleSize   string
	PreferredDate string
	Notes         string
}

func (s *bookingService) CreateBooking(ctx context.Context, input CreateBookingInput) (*database.Booking, error) {
	parsedDate, err := time.Parse("2006-01-02", input.PreferredDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	refCode := fmt.Sprintf("DCD-%04d", rand.Intn(9000)+1000)

	booking := database.Booking{
		ReferenceCode: refCode,
		Status:        database.StatusPendingQuote,
		LocationType:  input.LocationType,
		CustomerName:  input.CustomerName,
		CustomerEmail: input.CustomerEmail,
		CustomerPhone: input.CustomerPhone,
		Postcode:      input.Postcode,
		VehicleMake:   input.VehicleMake,
		VehicleModel:  input.VehicleModel,
		VehicleSize:   input.VehicleSize,
		PreferredDate: parsedDate,
		EstimatedCost: 140.0,
		Notes:         input.Notes,
	}

	if err := s.db.GetDB(ctx).Create(&booking).Error; err != nil {
		return nil, fmt.Errorf("failed to save booking: %w", err)
	}

	return &booking, nil
}
