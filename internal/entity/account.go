package entity

import "gorm.io/gorm"

type SaveBalance struct {
	gorm.Model
	Current   float64     `json:"current"`
	Withdrawn []Withdrawn `gorm:"foreignKey:SaveBalanceID"`
	UserID    uint
}
