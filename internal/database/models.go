package database

import (
	"time"

	"gorm.io/gorm"
)


type ServicePackage struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	ParentID    *uint            `json:"parent_id,omitempty"`
	Slug        string           `gorm:"uniqueIndex;not null" json:"slug"`
	Name        string           `gorm:"not null" json:"name"`
	Category    string           `json:"category"`
	BasePrice   float64          `gorm:"type:decimal(10,2);not null" json:"base_price"`
	DurationMin int              `json:"duration_min"`
	Description string           `json:"description,omitempty"`
	Subtypes    []ServicePackage `gorm:"foreignKey:ParentID" json:"subtypes,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

type Reservation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Contact   string         `gorm:"not null" json:"contact"`
	Service   string         `gorm:"not null" json:"service"`
	Date      time.Time      `gorm:"not null" json:"date"`
	Message   string         `json:"message"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}