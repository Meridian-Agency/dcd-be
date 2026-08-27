package service

import (
	"context"
	"fmt"
	"time"

	"dcd-be/internal/database"
)

type ReservationService interface {
	CreateReservation(ctx context.Context, input CreateReservationInput) (*database.Reservation, error)
	ListReservationsPublic(ctx context.Context) ([]PublicReservationDTO, error)
	ListReservationsAdmin(ctx context.Context) ([]database.Reservation, error)
}

type reservationService struct {
	db database.Service
}

func NewReservationService(db database.Service) ReservationService {
	return &reservationService{db: db}
}

type CreateReservationInput struct {
	Name    string
	Contact string
	Service string
	Date    string
	Message string
}

type PublicReservationDTO struct {
	Date time.Time `json:"date"`
}

func (s *reservationService) CreateReservation(ctx context.Context, input CreateReservationInput) (*database.Reservation, error) {
	parsedDate, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	reservation := database.Reservation{
		Name:    input.Name,
		Contact: input.Contact,
		Service: input.Service,
		Date:    parsedDate,
		Message: input.Message,
	}

	if err := s.db.GetDB(ctx).Create(&reservation).Error; err != nil {
		return nil, fmt.Errorf("failed to save reservation: %w", err)
	}

	return &reservation, nil
}

func (s *reservationService) ListReservationsPublic(ctx context.Context) ([]PublicReservationDTO, error) {
	var reservations []database.Reservation
	err := s.db.GetDB(ctx).
		Select("date").
		Order("date ASC").
		Find(&reservations).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to list public reservations: %w", err)
	}

	publicDTOs := make([]PublicReservationDTO, len(reservations))
	for i, r := range reservations {
		publicDTOs[i] = PublicReservationDTO{
			Date: r.Date,
		}
	}

	return publicDTOs, nil
}

func (s *reservationService) ListReservationsAdmin(ctx context.Context) ([]database.Reservation, error) {
	var reservations []database.Reservation
	err := s.db.GetDB(ctx).
		Order("date ASC").
		Find(&reservations).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to list admin reservations: %w", err)
	}

	return reservations, nil
}
