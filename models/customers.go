package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model

	Id        string    `gorm:"primaryKey"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Address   string    `gorm:"not null"`
	StartDate time.Time `gorm:"default:null"`
	EndDate   time.Time `gorm:"default:null"`
	IsActive  bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

const DeactivateCustomer string = "Order to deactivate customer"
