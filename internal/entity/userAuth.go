package entity

import (
	"github.com/golang-jwt/jwt/v5"
)

type ClaimUser struct {
	jwt.RegisteredClaims
	UserID uint
}
