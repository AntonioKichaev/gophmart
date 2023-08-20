package service

import (
	"context"
	"errors"

	"github.com/AntonioKichaev/internal/entity"
	"github.com/AntonioKichaev/internal/repository"
)

type AuthRepository interface {
	Register(ctx context.Context, dto *repository.RegisterDTO) (*entity.User, error)
	GetUserByLogin(ctx context.Context, login string) (*entity.User, error)
}
type AccountServiceAdapter interface {
	CreateBalance(ctx context.Context, dto CreateBalanceDTO) error
}

type Hasher interface {
	MakeHash(s string) string
}

type AuthService struct {
	r  AuthRepository
	h  Hasher
	as AccountServiceAdapter
}
type RegisterDTO struct {
	Login    string
	Password string
}
type LoginDTO struct {
	Login    string
	Password string
}

type hashe struct {
}

func (h hashe) MakeHash(s string) string {
	return s
}
func NewAuthService(r AuthRepository, as AccountServiceAdapter) *AuthService {
	return &AuthService{
		r:  r,
		h:  hashe{},
		as: as,
	}
}
func (us AuthService) Register(ctx context.Context, dto *RegisterDTO) (*entity.User, error) {
	// make hash
	user, err := us.r.Register(ctx, &repository.RegisterDTO{
		Login:    dto.Login,
		Password: us.h.MakeHash(dto.Password),
	})
	if err != nil {
		return nil, err
	}
	err = us.as.CreateBalance(ctx, CreateBalanceDTO{
		UserID: user.ID,
	})

	return user, err
}

func (us AuthService) Login(ctx context.Context, dto *LoginDTO) (*entity.User, error) {
	passhash := us.h.MakeHash(dto.Password)
	u, err := us.r.GetUserByLogin(ctx, dto.Login)
	if err != nil {
		return nil, errors.New("not found")
	}
	if u.Password != passhash {
		return nil, errors.New("errors password or login")
	}
	return u, err
}
