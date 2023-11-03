package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model

	Id        string `gorm:"primaryKey"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Address   string `gorm:"not null"`
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}
