package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"dcd-be/internal/config"
	"dcd-be/internal/database"
	"dcd-be/internal/service"

	"github.com/gin-gonic/gin"
)

type mockReservationService struct {
	createReservationFunc      func(ctx context.Context, input service.CreateReservationInput) (*database.Reservation, error)
	listReservationsPublicFunc func(ctx context.Context) ([]service.PublicReservationDTO, error)
	listReservationsAdminFunc  func(ctx context.Context) ([]database.Reservation, error)
}

func (m *mockReservationService) CreateReservation(ctx context.Context, input service.CreateReservationInput) (*database.Reservation, error) {
	return m.createReservationFunc(ctx, input)
}

func (m *mockReservationService) ListReservationsPublic(ctx context.Context) ([]service.PublicReservationDTO, error) {
	return m.listReservationsPublicFunc(ctx)
}

func (m *mockReservationService) ListReservationsAdmin(ctx context.Context) ([]database.Reservation, error) {
	return m.listReservationsAdminFunc(ctx)
}

func TestCreateReservationHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           string
		mockSetup      func() *mockReservationService
		expectedStatus int
		verifyBody     func(t *testing.T, body string)
	}{
		{
			name: "successful creation",
			body: `{"name":"John Doe","contact":"john@example.com","service":"Deep Clean Valet","date":"2026-09-01T10:00:00Z","message":"hello"}`,
			mockSetup: func() *mockReservationService {
				return &mockReservationService{
					createReservationFunc: func(ctx context.Context, input service.CreateReservationInput) (*database.Reservation, error) {
						tTime, _ := time.Parse(time.RFC3339, input.Date)
						return &database.Reservation{
							ID:        42,
							Name:      input.Name,
							Contact:   input.Contact,
							Service:   input.Service,
							Date:      tTime,
							Message:   input.Message,
							CreatedAt: time.Now(),
						}, nil
					},
				}
			},
			expectedStatus: http.StatusCreated,
			verifyBody: func(t *testing.T, body string) {
				var res database.Reservation
				if err := json.Unmarshal([]byte(body), &res); err != nil {
					t.Fatalf("failed to unmarshal body: %v", err)
				}
				if res.ID != 42 {
					t.Errorf("expected ID 42, got %d", res.ID)
				}
				if res.Name != "John Doe" {
					t.Errorf("expected Name 'John Doe', got '%s'", res.Name)
				}
			},
		},
		{
			name:           "invalid json / binding error",
			body:           `{"name":"","contact":""}`,
			mockSetup:      func() *mockReservationService { return &mockReservationService{} },
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body string) {
				if !strings.Contains(body, "error") {
					t.Errorf("expected error message, got %s", body)
				}
			},
		},
		{
			name: "service layer error",
			body: `{"name":"John Doe","contact":"john@example.com","service":"Deep Clean Valet","date":"2026-09-01T10:00:00Z"}`,
			mockSetup: func() *mockReservationService {
				return &mockReservationService{
					createReservationFunc: func(ctx context.Context, input service.CreateReservationInput) (*database.Reservation, error) {
						return nil, errors.New("db error")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyBody: func(t *testing.T, body string) {
				if !strings.Contains(body, "Failed to store reservation") {
					t.Errorf("expected error message, got %s", body)
				}
			},
		},
		{
			name: "duplicate reservation conflict",
			body: `{"name":"John Doe","contact":"john@example.com","service":"Deep Clean Valet","date":"2026-09-01T10:00:00Z"}`,
			mockSetup: func() *mockReservationService {
				return &mockReservationService{
					createReservationFunc: func(ctx context.Context, input service.CreateReservationInput) (*database.Reservation, error) {
						return nil, service.ErrDuplicateReservation
					},
				}
			},
			expectedStatus: http.StatusConflict,
			verifyBody: func(t *testing.T, body string) {
				if !strings.Contains(body, service.ErrDuplicateReservation.Error()) {
					t.Errorf("expected duplicate error, got %s", body)
				}
			},
		},
		{
			name: "timeslot taken conflict",
			body: `{"name":"John Doe","contact":"john@example.com","service":"Deep Clean Valet","date":"2026-09-01T10:00:00Z"}`,
			mockSetup: func() *mockReservationService {
				return &mockReservationService{
					createReservationFunc: func(ctx context.Context, input service.CreateReservationInput) (*database.Reservation, error) {
						return nil, service.ErrSlotTaken
					},
				}
			},
			expectedStatus: http.StatusConflict,
			verifyBody: func(t *testing.T, body string) {
				if !strings.Contains(body, service.ErrSlotTaken.Error()) {
					t.Errorf("expected slot taken error, got %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				reservationService: tt.mockSetup(),
			}

			router := gin.New()
			router.Use(ErrorHandler())
			router.POST("/api/v1/reservations", s.createReservationHandler)

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/reservations", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			tt.verifyBody(t, w.Body.String())
		})
	}
}

func TestListReservationsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockReservationService{
		listReservationsPublicFunc: func(ctx context.Context) ([]service.PublicReservationDTO, error) {
			return []service.PublicReservationDTO{
				{Date: time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)},
			}, nil
		},
	}

	s := &Server{
		reservationService: mockSvc,
	}

	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/api/v1/reservations", s.listReservationsHandler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/reservations", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var publicList []service.PublicReservationDTO
	if err := json.Unmarshal(w.Body.Bytes(), &publicList); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(publicList) != 1 {
		t.Errorf("expected 1 item, got %d", len(publicList))
	}
}

func TestAdminReservationsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockReservationService{
		listReservationsAdminFunc: func(ctx context.Context) ([]database.Reservation, error) {
			return []database.Reservation{
				{ID: 1, Name: "Admin Test", Contact: "admin@test.com", Service: "S1", Date: time.Now()},
			}, nil
		},
	}

	s := &Server{
		cfg: &config.Config{
			AdminAPIKey: "secret-token",
		},
		reservationService: mockSvc,
	}

	router := gin.New()
	router.Use(ErrorHandler())

	admin := router.Group("/api/v1/admin")
	admin.Use(s.adminAuthMiddleware())
	{
		admin.GET("/reservations", s.listAdminReservationsHandler)
	}

	t.Run("Unauthorized - No Token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/reservations", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})

	t.Run("Unauthorized - Bad Token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/reservations", nil)
		req.Header.Set("Authorization", "Bearer bad-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})

	t.Run("Authorized - Success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/reservations", nil)
		req.Header.Set("Authorization", "Bearer secret-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var list []database.Reservation
		if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if len(list) != 1 || list[0].Name != "Admin Test" {
			t.Errorf("expected list with Admin Test, got %+v", list)
		}
	})
}
