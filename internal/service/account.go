package service

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/AntonioKichaev/internal/entity"
	"github.com/AntonioKichaev/internal/model"
	"github.com/AntonioKichaev/internal/repository"
)

//go:generate mockery --name=AccountRepository --output=../mocks --filename=account_repository.go
type AccountRepository interface {
	GetBalanceByID(ctx context.Context, dto repository.GetBalanceByIDDTO) (*entity.SaveBalance, error)
	GetWithdrawsByUserID(ctx context.Context, dto repository.GetWithdrawsByUserIDDTO) ([]entity.Withdrawn, error)
	UpdateAccount(ctx context.Context, balance entity.SaveBalance) error
	CreateBalance(ctx context.Context, dto repository.CreateBalanceDTO) error
	AddPoint(ctx context.Context, dto repository.AddPointDTO) error
}

//go:generate mockery --name=OrderAdapter --output=../mocks --filename=order_adapter.go
type OrderAdapter interface {
	GetOrderByNumber(ctx context.Context, dto GetOrderByNumberDTO) (*entity.Order, error)
	UploadOrderID(ctx context.Context, dto *UploadOrderIDDTO) (*entity.Order, error)
	GetOrdersByUser(ctx context.Context, dto *GetOrdersByUserIDDTO) ([]entity.Order, error)
}

//go:generate mockery --name=WithdrawnAdapter --output=../mocks --filename=withdrawn_adapter.go
type WithdrawnAdapter interface {
	Create(ctx context.Context, dto WithdrawnCreateDTO) error
	GetWithdraws(ctx context.Context, userID uint) ([]entity.Withdrawn, error)
}
type Account struct {
	r  AccountRepository
	oa OrderAdapter
	wd WithdrawnAdapter
}

func NewAccountService(r AccountRepository, oa OrderAdapter, wd WithdrawnAdapter) *Account {
	return &Account{
		r:  r,
		oa: oa,
		wd: wd,
	}
}

type GetBalanceByIDDTO struct {
	UserID uint
}

type AddPointDTO struct {
	UserID uint
	Point  float64
}
type WithdrawByUserIDDTO struct {
	UserID  uint
	OrderID string
	Amount  float64
}
type WithdrawsByUserIDDTO struct {
	UserID uint
}

type CreateBalanceDTO struct {
	UserID uint
}

func (a *Account) GetBalanceByID(ctx context.Context, dto *GetBalanceByIDDTO) (*model.UserBalance, error) {

	sb, err := a.r.GetBalanceByID(ctx, repository.GetBalanceByIDDTO{
		UserID: dto.UserID,
	})

	if err != nil {
		return nil, err
	}
	wd, err := a.r.GetWithdrawsByUserID(ctx, repository.GetWithdrawsByUserIDDTO{
		UserID: dto.UserID,
	})

	if err != nil {
		return nil, err
	}
	wdTotalSpent := 0.0
	for _, w := range wd {
		wdTotalSpent += w.Sum
	}
	ub := &model.UserBalance{
		Current:   sb.Current,
		Withdrawn: wdTotalSpent,
	}

	return ub, nil
}

var ErrorNotEnoughMoney = errors.New("not enough money")
var ErrorNotFoundOrder = errors.New("not found order")

func (a *Account) Withdrawn(ctx context.Context, dto *WithdrawByUserIDDTO) error {
	sb, err := a.r.GetBalanceByID(ctx, repository.GetBalanceByIDDTO{
		UserID: dto.UserID,
	})
	if err != nil {
		return err
	}

	if sb.Current < dto.Amount {
		return ErrorNotEnoughMoney
	}

	order, err := a.oa.UploadOrderID(ctx, &UploadOrderIDDTO{
		Number: dto.OrderID,
		UserID: dto.UserID,
	})

	if err != nil {
		return err
	}

	sb.Current = sb.Current - dto.Amount
	err = a.r.UpdateAccount(ctx, *sb)
	if err != nil {
		return err
	}

	err = a.wd.Create(ctx, WithdrawnCreateDTO{
		OrderID:       order.ID,
		SaveBalanceID: sb.ID,
		UserID:        dto.UserID,
		Sum:           dto.Amount,
	})

	return err
}
func (a *Account) GetWithdraws(ctx context.Context, dto *WithdrawsByUserIDDTO) ([]model.Withdrawal, error) {
	withdrawns, err := a.wd.GetWithdraws(ctx, dto.UserID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	orders, err := a.oa.GetOrdersByUser(ctx, &GetOrdersByUserIDDTO{
		UserID: dto.UserID,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	uniqOrder := make(map[uint]entity.Order, len(orders))

	for _, order := range orders {
		uniqOrder[order.ID] = order
	}
	var result []model.Withdrawal

	for _, w := range withdrawns {
		o, ok := uniqOrder[w.OrderID]
		if !ok {
			continue
		}
		result = append(result, model.Withdrawal{
			Order:       o.Number,
			Sum:         w.Sum,
			ProcessedAt: w.ProcessedAt,
		})
	}

	return result, nil
}

func (a *Account) CreateBalance(ctx context.Context, dto CreateBalanceDTO) error {
	err := a.r.CreateBalance(ctx, repository.CreateBalanceDTO{
		UserID: dto.UserID,
	})
	return err
}

func (a *Account) AddPoint(ctx context.Context, dto AddPointDTO) error {
	return a.r.AddPoint(ctx, repository.AddPointDTO{
		UserID: dto.UserID,
		Point:  dto.Point,
	})
}
