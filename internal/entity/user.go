package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login       string      `json:"login" gorm:"index:idx_login,unique"`
	Password    string      `json:"password"`
	Order       []Order     `gorm:"foreignKey:UserID"`
	SaveBalance SaveBalance `gorm:"foreignKey:UserID"`
}
