package server

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// Centralized Error Handling Middleware
	r.Use(ErrorHandler())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://demon-car-detailing-phi.vercel.app", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", s.healthHandler)
		v1.GET("/ping", s.pingHandler)
		v1.GET("/coverage/check", s.checkCoverageHandler)
		v1.POST("/bookings", s.createBookingHandler)
	}

	return r
}
