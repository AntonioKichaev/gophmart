package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AntonioKichaev/internal/entity"
	"github.com/AntonioKichaev/internal/service"
)

type RequestRegister struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RequestLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserService interface {
	Register(ctx context.Context, dto *service.RegisterDTO) (*entity.User, error)
	Login(ctx context.Context, dto *service.LoginDTO) (*entity.User, error)
}
type AuthAdapter interface {
	BuildJWTString(userID uint) (string, error)
}
type Auth struct {
	us UserService
	aa AuthAdapter
}

func NewAuthHandle(us UserService, aa AuthAdapter) *Auth {
	return &Auth{
		us: us,
		aa: aa,
	}
}

func (h *Auth) Register(c *gin.Context) {
	r := &RequestRegister{}
	err := c.BindJSON(r)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if !isValidRequestRegister(r) {
		c.AbortWithStatus(http.StatusBadRequest)
		//400 - неверный формат
		return
	}
	user, err := h.us.Register(c.Request.Context(), &service.RegisterDTO{
		Login:    r.Login,
		Password: r.Password,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusConflict)
		return
		//409 - логин занят
		//500 - ошибка
	}
	token, err := h.aa.BuildJWTString(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, `{"error": "buildJWTString"}`)
		return
	}
	c.Writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	c.JSON(http.StatusOK, user)

}

func (h *Auth) Login(c *gin.Context) {
	r := &RequestLogin{}
	err := c.BindJSON(r)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if !isValidRequestLogin(r) {
		//400 - неверный формат
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	u, err := h.us.Login(c.Request.Context(), &service.LoginDTO{
		Login:    r.Login,
		Password: r.Password,
	})
	if err != nil {
		//401 - неверная пара логин/пароль
		//500 - ошибка
		c.JSON(http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	token, err := h.aa.BuildJWTString(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, `{"error": "buildJWTString"}`)
		return
	}
	c.Writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	c.JSON(http.StatusOK, r)
}

func isValidRequestRegister(r *RequestRegister) bool {
	if len(r.Login) == 0 {
		return false
	}
	if len(r.Password) == 0 {
		return false
	}
	return true
}

func isValidRequestLogin(r *RequestLogin) bool {
	if len(r.Login) == 0 {
		return false
	}
	if len(r.Password) == 0 {
		return false
	}
	return true
}
