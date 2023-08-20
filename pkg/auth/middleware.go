package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/AntonioKichaev/internal/entity"
)

type AuthMid struct {
	SecretKey string
}

func NewAuthMid(secretKey string) *AuthMid {
	return &AuthMid{
		SecretKey: secretKey,
	}
}
func AuthMiddleware(a *AuthMid) func(c *gin.Context) {

	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		auth = strings.TrimPrefix(auth, "Bearer ")
		if len(auth) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		currentUser, err := a.GetUserID(auth)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("UserID", currentUser)
	}

}

func (a *AuthMid) BuildJWTString(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.ClaimUser{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(a.SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func (a *AuthMid) GetUserID(tokenString string) (uint, error) {
	claims := &entity.ClaimUser{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.SecretKey), nil
	})

	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, err
	}

	// возвращаем UserID пользователя в читаемом виде
	return claims.UserID, nil
}
