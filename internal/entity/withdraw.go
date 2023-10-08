package entity

import (
	"time"

	"gorm.io/gorm"
)

type Withdrawn struct {
	gorm.Model
	OrderID       uint
	SaveBalanceID uint
	UserID        uint
	Sum           float64   `json:"sum"`
	ProcessedAt   time.Time `json:"processed_at"`
}
