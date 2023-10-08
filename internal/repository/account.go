package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/AntonioKichaev/internal/entity"
)

type AccountRepository struct {
	db *gorm.DB
}

type GetBalanceByIDDTO struct {
	UserID uint
}

type GetWithdrawsByUserIDDTO struct {
	UserID uint
}
type CreateBalanceDTO struct {
	UserID uint
}
type AddPointDTO struct {
	UserID uint
	Point  float64
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}
func (a *AccountRepository) GetBalanceByID(ctx context.Context, dto GetBalanceByIDDTO) (*entity.SaveBalance, error) {
	sb := entity.SaveBalance{}
	tx := a.db.Model(entity.SaveBalance{}).Where("user_id = ?", dto.UserID).First(&sb)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &sb, nil
}

func (a *AccountRepository) GetWithdrawsByUserID(ctx context.Context, dto GetWithdrawsByUserIDDTO) ([]entity.Withdrawn, error) {
	var wd []entity.Withdrawn
	tx := a.db.Model(entity.Withdrawn{}).Where("user_id = ?", dto.UserID).Find(&wd)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return wd, nil
}

func (a *AccountRepository) UpdateAccount(ctx context.Context, balance entity.SaveBalance) error {

	tr := a.db.Model(entity.SaveBalance{}).Where("id = ?", balance.ID).Updates(balance)
	return tr.Error
}

func (a *AccountRepository) CreateBalance(ctx context.Context, dto CreateBalanceDTO) error {
	sb := entity.SaveBalance{
		UserID:  dto.UserID,
		Current: 0.0,
	}
	tr := a.db.Model(entity.SaveBalance{}).Create(&sb)
	return tr.Error
}

func (a *AccountRepository) AddPoint(ctx context.Context, dto AddPointDTO) error {
	sb := &entity.SaveBalance{
		UserID: dto.UserID,
	}
	tx := a.db.Model(entity.SaveBalance{}).Where("user_id = ?", dto.UserID).First(&sb)
	if tx.Error != nil {
		return tx.Error
	}
	sb.Current += dto.Point
	tx = a.db.Model(entity.SaveBalance{}).Where("id=?", sb.ID).Updates(sb)
	return tx.Error
}
