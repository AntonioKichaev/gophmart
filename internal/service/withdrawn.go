package service

import (
	"context"

	"github.com/AntonioKichaev/internal/entity"
	"github.com/AntonioKichaev/internal/repository"
)

type WithdrawnRepository interface {
	Create(ctx context.Context, dto repository.CreateWithdrawnDTO) error
	GetAllByUserID(ctx context.Context, userID uint) ([]entity.Withdrawn, error)
}

type WithdrawnService struct {
	r WithdrawnRepository
}

func NewWithdrawnService(r WithdrawnRepository) *WithdrawnService {
	return &WithdrawnService{
		r: r,
	}
}

type WithdrawnCreateDTO struct {
	OrderID       uint
	SaveBalanceID uint
	UserID        uint
	Sum           float64
}

func (w *WithdrawnService) Create(ctx context.Context, dto WithdrawnCreateDTO) error {
	err := w.r.Create(ctx, repository.CreateWithdrawnDTO{
		OrderID:       dto.OrderID,
		SaveBalanceID: dto.SaveBalanceID,
		UserID:        dto.UserID,
		Sum:           dto.Sum,
	})
	return err
}

func (w *WithdrawnService) GetWithdraws(ctx context.Context, userID uint) ([]entity.Withdrawn, error) {
	return w.r.GetAllByUserID(ctx, userID)
}
