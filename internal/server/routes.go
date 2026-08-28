package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set("RequestID", reqID)
		c.Header("X-Request-ID", reqID)
		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		reqIDVal, _ := c.Get("RequestID")
		reqID, _ := reqIDVal.(string)

		slog.InfoContext(c.Request.Context(), "HTTP request processed",
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Int("status", c.Writer.Status()),
			slog.String("ip", c.ClientIP()),
			slog.Duration("latency", latency),
			slog.String("request_id", reqID),
		)
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.New()

	r.Use(RequestIDMiddleware())
	r.Use(LoggerMiddleware())
	r.Use(gin.Recovery())
	r.Use(ErrorHandler())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     s.cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", s.healthHandler)
		v1.GET("/ping", s.pingHandler)
		v1.GET("/services", s.listServicesHandler)
		v1.GET("/services/:slug", s.getServiceBySlugHandler)

		// Reservation routes
		v1.POST("/reservations", s.createReservationHandler)
		v1.GET("/reservations", s.listReservationsHandler)

		// Admin reservation routes
		admin := v1.Group("/admin")
		admin.Use(s.adminAuthMiddleware())
		{
			admin.GET("/reservations", s.listAdminReservationsHandler)
		}
	}

	return r
}

func (s *Server) adminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		expectedToken := "Bearer " + s.cfg.AdminAPIKey
		if token != expectedToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}
