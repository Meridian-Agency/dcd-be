package database

import (
	"time"

	"gorm.io/gorm"
)

type BookingStatus string

const (
	StatusPendingQuote BookingStatus = "PENDING_QUOTE"
	StatusConfirmed    BookingStatus = "CONFIRMED"
	StatusCancelled    BookingStatus = "CANCELLED"
)

type ServicePackage struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Slug        string    `gorm:"uniqueIndex;not null" json:"slug"`
	Name        string    `gorm:"not null" json:"name"`
	Category    string    `json:"category"`
	BasePrice   float64   `gorm:"type:decimal(10,2);not null" json:"base_price"`
	DurationMin int       `json:"duration_min"`
	CreatedAt   time.Time `json:"created_at"`
}

type Booking struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ReferenceCode string         `gorm:"uniqueIndex;not null" json:"reference_code"`
	Status        BookingStatus  `gorm:"default:'PENDING_QUOTE'" json:"status"`
	LocationType  string         `gorm:"not null" json:"location_type"`
	CustomerName  string         `gorm:"not null" json:"customer_name"`
	CustomerEmail string         `gorm:"not null" json:"customer_email"`
	CustomerPhone string         `gorm:"not null" json:"customer_phone"`
	Postcode      string         `json:"postcode,omitempty"`
	VehicleMake   string         `gorm:"not null" json:"vehicle_make"`
	VehicleModel  string         `gorm:"not null" json:"vehicle_model"`
	VehicleSize   string         `gorm:"not null" json:"vehicle_size"`
	PreferredDate time.Time      `gorm:"not null" json:"preferred_date"`
	EstimatedCost float64        `gorm:"type:decimal(10,2)" json:"estimated_cost"`
	Notes         string         `json:"notes,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}