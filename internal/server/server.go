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
	coverageService       service.CoverageService
	servicePackageService service.ServicePackageService
}

func NewServer(cfg *config.Config) *http.Server {
	dbService := database.New(cfg)

	s := &Server{
		cfg:                   cfg,
		db:                    dbService,
		reservationService:    service.NewReservationService(dbService),
		coverageService:       service.NewCoverageService(),
		servicePackageService: service.NewServicePackageService(dbService),
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
