package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"dcd-be/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service interface {
	Health() map[string]string
	Close() error
	GetDB(ctx context.Context) *gorm.DB
	Transactor
}

type service struct {
	db *gorm.DB
}

var dbInstance *service

func New(cfg *config.Config) Service {
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBDatabase, cfg.DBSchema)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{db: db}
	return dbInstance
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&ServicePackage{}, &Reservation{})
}

func (s *service) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionKey).(*gorm.DB); ok {
		return tx
	}
	return s.db.WithContext(ctx)
}

func (s *service) WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, transactionKey, tx)
		return fn(txCtx)
	})
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sqlDB, err := s.db.DB()
	if err != nil {
		return map[string]string{"status": "down", "error": err.Error()}
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return map[string]string{"status": "down", "error": err.Error()}
	}

	return map[string]string{"status": "up", "message": "DCD database is healthy"}
}

func (s *service) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}