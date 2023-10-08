package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/AntonioKichaev/internal/entity"
)

type WithdrawnRepository struct {
	db *gorm.DB
}

type CreateWithdrawnDTO struct {
	OrderID       uint
	SaveBalanceID uint
	UserID        uint
	Sum           float64
}

func NewWithdrawnRepository(db *gorm.DB) *WithdrawnRepository {
	return &WithdrawnRepository{
		db: db,
	}
}

func (w *WithdrawnRepository) Create(ctx context.Context, dto CreateWithdrawnDTO) error {
	tx := w.db.Model(entity.Withdrawn{}).Create(&entity.Withdrawn{
		OrderID:       dto.OrderID,
		SaveBalanceID: dto.SaveBalanceID,
		UserID:        dto.UserID,
		Sum:           dto.Sum,
		ProcessedAt:   time.Now(),
	})
	return tx.Error
}

func (w *WithdrawnRepository) GetAllByUserID(ctx context.Context, userID uint) ([]entity.Withdrawn, error) {
	var wd []entity.Withdrawn
	tx := w.db.Model(entity.Withdrawn{}).Where("user_id = ?", userID).Find(&wd)
	return wd, tx.Error
}
