package entity

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model `json:"-"`
	Number     string      `json:"number"`
	Status     string      `json:"status"`
	Accrual    float64     `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
	UserID     uint        `json:"-"`
	Withdrawn  []Withdrawn `gorm:"foreignKey:OrderID" json:"-"`
}
