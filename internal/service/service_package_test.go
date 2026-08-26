package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"dcd-be/internal/database"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbHost string
	dbPort int
	dbName string
	dbUser string
	dbPass string
)

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	dbName = "dcd_db_test"
	dbPass = "password"
	dbUser = "user"

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	host, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	port, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	portInt, err := strconv.Atoi(port.Port())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbHost = host
	dbPort = portInt

	return dbContainer.Terminate, nil
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	m.Run()

	if teardown != nil {
		_ = teardown(context.Background())
	}
}

type testDBService struct {
	db *gorm.DB
}

func (t *testDBService) Health() map[string]string { return nil }
func (t *testDBService) Close() error               { return nil }
func (t *testDBService) GetDB(ctx context.Context) *gorm.DB {
	return t.db.WithContext(ctx)
}
func (t *testDBService) WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return fn(ctx)
}

func setupTestDB(t *testing.T) database.Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(postgresDriver.Open(connStr), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open postgres test db: %v", err)
	}

	// Clean up / truncate tables if they exist
	db.Exec("DROP TABLE IF EXISTS bookings")
	db.Exec("DROP TABLE IF EXISTS service_packages")

	err = database.Migrate(db)
	if err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}

	return &testDBService{db: db}
}

func TestServicePackageService_ListServices(t *testing.T) {
	dbSrv := setupTestDB(t)
	db := dbSrv.GetDB(context.Background())

	// Run seeder to populate DB
	err := database.SeedServices(db)
	if err != nil {
		t.Fatalf("failed to seed services: %v", err)
	}

	svc := NewServicePackageService(dbSrv)
	ctx := context.Background()

	services, err := svc.ListServices(ctx)
	if err != nil {
		t.Fatalf("ListServices returned error: %v", err)
	}

	// Verify that we only get top-level packages (where parent_id is nil)
	// We expect 9 top-level packages:
	// 1. deep-clean-valet
	// 2. maintenance-valet
	// 3. single-stage-enhancement
	// 4. two-stage-paint-correction
	// 5. gtechniq-ceramic-coating
	// 6. paint-protection-film
	// 7. alloy-wheel-refurbishment
	// 8. vehicle-upgrades
	// 9. apple-carplay
	if len(services) != 9 {
		t.Errorf("expected 9 top-level services, got %d", len(services))
	}

	// Verify that paint-protection-film has preloaded subtypes
	var ppf *database.ServicePackage
	for i := range services {
		if services[i].Slug == "paint-protection-film" {
			ppf = &services[i]
			break
		}
	}

	if ppf == nil {
		t.Fatal("paint-protection-film package not found in list")
	}

	if len(ppf.Subtypes) != 3 {
		t.Errorf("expected 3 subtypes for PPF, got %d", len(ppf.Subtypes))
	}

	// Verify the subtypes exist
	expectedSubtypes := map[string]bool{
		"ppf-front-end": true,
		"ppf-extended":  true,
		"ppf-full-body": true,
	}

	for _, sub := range ppf.Subtypes {
		if !expectedSubtypes[sub.Slug] {
			t.Errorf("unexpected PPF subtype: %s", sub.Slug)
		}
		if *sub.ParentID != ppf.ID {
			t.Errorf("expected subtype ParentID to be %d, got %d", ppf.ID, *sub.ParentID)
		}
	}
}

func TestServicePackageService_GetServiceBySlug(t *testing.T) {
	dbSrv := setupTestDB(t)
	db := dbSrv.GetDB(context.Background())

	err := database.SeedServices(db)
	if err != nil {
		t.Fatalf("failed to seed services: %v", err)
	}

	svc := NewServicePackageService(dbSrv)
	ctx := context.Background()

	// Get a specific top-level service
	service, err := svc.GetServiceBySlug(ctx, "paint-protection-film")
	if err != nil {
		t.Fatalf("failed to get paint-protection-film: %v", err)
	}

	if service.Name != "Paint Protection Film (PPF)" {
		t.Errorf("expected name 'Paint Protection Film (PPF)', got %s", service.Name)
	}

	if len(service.Subtypes) != 3 {
		t.Errorf("expected preloaded subtypes, got %d", len(service.Subtypes))
	}

	// Get a subtype directly
	subtype, err := svc.GetServiceBySlug(ctx, "ppf-front-end")
	if err != nil {
		t.Fatalf("failed to get ppf-front-end: %v", err)
	}

	if subtype.Name != "Front End Coverage" {
		t.Errorf("expected name 'Front End Coverage', got %s", subtype.Name)
	}

	if subtype.ParentID == nil {
		t.Fatal("expected ParentID to be non-nil for subtype")
	}
}
