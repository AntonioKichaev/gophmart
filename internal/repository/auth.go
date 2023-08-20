package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/AntonioKichaev/internal/entity"
)

type AuthRepository struct {
	db *gorm.DB
}

type RegisterDTO struct {
	Login    string
	Password string
}

type LoginDTO struct {
	Login    string
	Password string
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}
func (a AuthRepository) Register(ctx context.Context, dto *RegisterDTO) (*entity.User, error) {
	user := &entity.User{
		Login:    dto.Login,
		Password: dto.Password,
	}
	tx := a.db.Create(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (a AuthRepository) GetUserByLogin(ctx context.Context, login string) (*entity.User, error) {
	u := entity.User{}
	tx := a.db.Model(entity.User{}).Where("login = ?", login).First(&u)
	return &u, tx.Error
}
