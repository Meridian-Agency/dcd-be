package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dcd-be/internal/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type mockServicePackageService struct {
	listServicesFunc     func(ctx context.Context) ([]database.ServicePackage, error)
	getServiceBySlugFunc func(ctx context.Context, slug string) (*database.ServicePackage, error)
}

func (m *mockServicePackageService) ListServices(ctx context.Context) ([]database.ServicePackage, error) {
	return m.listServicesFunc(ctx)
}

func (m *mockServicePackageService) GetServiceBySlug(ctx context.Context, slug string) (*database.ServicePackage, error) {
	return m.getServiceBySlugFunc(ctx, slug)
}

func TestListServicesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func() *mockServicePackageService
		expectedStatus int
		verifyBody     func(t *testing.T, body string)
	}{
		{
			name: "successful retrieval",
			mockSetup: func() *mockServicePackageService {
				return &mockServicePackageService{
					listServicesFunc: func(ctx context.Context) ([]database.ServicePackage, error) {
						return []database.ServicePackage{
							{
								ID:          1,
								Slug:        "test-service",
								Name:        "Test Service",
								Category:    "VALETING",
								BasePrice:   100.0,
								DurationMin: 60,
								Description: "Test description",
								CreatedAt:   time.Now(),
							},
						}, nil
					},
				}
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body string) {
				var services []database.ServicePackage
				if err := json.Unmarshal([]byte(body), &services); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if len(services) != 1 {
					t.Fatalf("expected 1 service, got %d", len(services))
				}
				if services[0].Slug != "test-service" {
					t.Errorf("expected slug 'test-service', got '%s'", services[0].Slug)
				}
			},
		},
		{
			name: "service layer error",
			mockSetup: func() *mockServicePackageService {
				return &mockServicePackageService{
					listServicesFunc: func(ctx context.Context) ([]database.ServicePackage, error) {
						return nil, errors.New("db error")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyBody: func(t *testing.T, body string) {
				var resp map[string]string
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp["error"] != "Failed to retrieve services" {
					t.Errorf("expected error message 'Failed to retrieve services', got '%s'", resp["error"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				servicePackageService: tt.mockSetup(),
			}

			router := gin.New()
			router.Use(ErrorHandler())
			router.GET("/api/v1/services", s.listServicesHandler)

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/services", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			tt.verifyBody(t, w.Body.String())
		})
	}
}

func TestGetServiceBySlugHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		slug           string
		mockSetup      func() *mockServicePackageService
		expectedStatus int
		verifyBody     func(t *testing.T, body string)
	}{
		{
			name: "successful retrieval",
			slug: "test-service",
			mockSetup: func() *mockServicePackageService {
				return &mockServicePackageService{
					getServiceBySlugFunc: func(ctx context.Context, slug string) (*database.ServicePackage, error) {
						return &database.ServicePackage{
							ID:          1,
							Slug:        slug,
							Name:        "Test Service",
							Category:    "VALETING",
							BasePrice:   100.0,
							DurationMin: 60,
							Description: "Test description",
							CreatedAt:   time.Now(),
						}, nil
					},
				}
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body string) {
				var service database.ServicePackage
				if err := json.Unmarshal([]byte(body), &service); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if service.Slug != "test-service" {
					t.Errorf("expected slug 'test-service', got '%s'", service.Slug)
				}
			},
		},
		{
			name: "not found error",
			slug: "non-existent",
			mockSetup: func() *mockServicePackageService {
				return &mockServicePackageService{
					getServiceBySlugFunc: func(ctx context.Context, slug string) (*database.ServicePackage, error) {
						return nil, gorm.ErrRecordNotFound
					},
				}
			},
			expectedStatus: http.StatusNotFound,
			verifyBody: func(t *testing.T, body string) {
				var resp map[string]string
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp["error"] != "Service not found" {
					t.Errorf("expected error message 'Service not found', got '%s'", resp["error"])
				}
			},
		},
		{
			name: "service layer error",
			slug: "test-service",
			mockSetup: func() *mockServicePackageService {
				return &mockServicePackageService{
					getServiceBySlugFunc: func(ctx context.Context, slug string) (*database.ServicePackage, error) {
						return nil, errors.New("db error")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyBody: func(t *testing.T, body string) {
				var resp map[string]string
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp["error"] != "Failed to retrieve service" {
					t.Errorf("expected error message 'Failed to retrieve service', got '%s'", resp["error"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				servicePackageService: tt.mockSetup(),
			}

			router := gin.New()
			router.Use(ErrorHandler())
			router.GET("/api/v1/services/:slug", s.getServiceBySlugHandler)

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/services/"+tt.slug, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			tt.verifyBody(t, w.Body.String())
		})
	}
}
