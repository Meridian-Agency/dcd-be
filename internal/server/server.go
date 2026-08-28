package server

import (
	"fmt"
	"net/http"
	"time"

	"dcd-be/internal/config"
	"dcd-be/internal/database"
	"dcd-be/internal/service"
)

type Server struct {
	cfg                   *config.Config
	db                    database.Service
	reservationService    service.ReservationService
	servicePackageService service.ServicePackageService
}

func NewHandler(cfg *config.Config, dbService database.Service) http.Handler {
	s := &Server{
		cfg:                   cfg,
		db:                    dbService,
		reservationService:    service.NewReservationService(dbService),
		servicePackageService: service.NewServicePackageService(dbService),
	}

	return s.RegisterRoutes()
}

func NewServer(cfg *config.Config, dbService database.Service) *http.Server {
	handler := NewHandler(cfg, dbService)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
