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
	cfg             *config.Config
	db              database.Service
	bookingService  service.BookingService
	coverageService service.CoverageService
}

func NewServer(cfg *config.Config) *http.Server {
	dbService := database.New(cfg)

	s := &Server{
		cfg:             cfg,
		db:              dbService,
		bookingService:  service.NewBookingService(dbService),
		coverageService: service.NewCoverageService(),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
