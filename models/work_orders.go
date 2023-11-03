package models

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	New       Status = "new"
	Done      Status = "done"
	Cancelled Status = "cancelled"
)

type WorkOrder struct {
	gorm.Model

	Id               string `json:"id" gorm:"primarykey"`
	CustomerId       string `json:"customer_id" gorm:"references:Id;not null"`
	Customer         Customer
	Title            string    `json:"title"`
	PlannedDateBegin time.Time `json:"planned_date_begin"`
	PlannedDateEnd   time.Time `json:"planned_date_end"`
	Status           Status    `json:"status"`
}
