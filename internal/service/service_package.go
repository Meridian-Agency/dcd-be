package service

import (
	"context"
	"fmt"

	"dcd-be/internal/database"
)

type ServicePackageService interface {
	ListServices(ctx context.Context) ([]database.ServicePackage, error)
	GetServiceBySlug(ctx context.Context, slug string) (*database.ServicePackage, error)
}

type servicePackageService struct {
	db database.Service
}

func NewServicePackageService(db database.Service) ServicePackageService {
	return &servicePackageService{db: db}
}

func (s *servicePackageService) ListServices(ctx context.Context) ([]database.ServicePackage, error) {
	var services []database.ServicePackage
	err := s.db.GetDB(ctx).
		Where("parent_id IS NULL").
		Preload("Subtypes").
		Order("id ASC").
		Find(&services).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to list service packages: %w", err)
	}

	return services, nil
}

func (s *servicePackageService) GetServiceBySlug(ctx context.Context, slug string) (*database.ServicePackage, error) {
	var service database.ServicePackage
	err := s.db.GetDB(ctx).
		Where("slug = ?", slug).
		Preload("Subtypes").
		First(&service).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to get service package by slug %s: %w", slug, err)
	}

	return &service, nil
}
