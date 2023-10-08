package service

import (
	"context"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/AntonioKichaev/internal/entity"
	"github.com/AntonioKichaev/internal/repository"
)

type OrderRepository interface {
	Create(ctx context.Context, dto repository.CreateOrderDTO) (entity.Order, error)
	GetOrders(ctx context.Context, dto repository.GetOrdersDTO) []entity.Order
	GetOrder(ctx context.Context, dto repository.GetOrderByNumberDTO) (*entity.Order, error)
	UpdateOrder(ctx context.Context, dto repository.UpdateOrderDTO) error
	GetOrdersInProgress(ctx context.Context) ([]entity.Order, error)
}

type Order struct {
	r OrderRepository
}

func NewOrder(r OrderRepository) *Order {
	return &Order{
		r: r,
	}
}

type UploadOrderIDDTO struct {
	UserID uint
	Number string
}

type UpdateOrderDTO struct {
	Number  string
	Status  string
	Accrual float64
	OrderID uint
}
type GetOrdersByUserIDDTO struct {
	UserID uint
}

type GetOrderByNumberDTO struct {
	Number string
}

var ErrOrderExists = errors.New("founded")
var ErrOrderUploadedOtherUser = errors.New("other user")

func (o *Order) UploadOrderID(ctx context.Context, dto *UploadOrderIDDTO) (*entity.Order, error) {
	order, err := o.r.GetOrder(ctx, repository.GetOrderByNumberDTO{
		Number: dto.Number,
	})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else {
		if order.UserID == dto.UserID {
			return order, ErrOrderExists
		}
		return nil, ErrOrderUploadedOtherUser

	}

	orderN, err := o.r.Create(ctx, repository.CreateOrderDTO{
		UserID:     dto.UserID,
		Number:     dto.Number,
		Status:     "NEW",
		UploadedAt: time.Now(),
	})
	return &orderN, err
}

func (o *Order) GetOrdersByUser(ctx context.Context, dto *GetOrdersByUserIDDTO) ([]entity.Order, error) {
	orders := o.r.GetOrders(ctx, repository.GetOrdersDTO{
		UserID: dto.UserID,
	})
	return orders, nil
}

func (o *Order) GetOrderByNumber(ctx context.Context, dto GetOrderByNumberDTO) (*entity.Order, error) {
	order, err := o.r.GetOrder(ctx, repository.GetOrderByNumberDTO{
		Number: dto.Number,
	})
	return order, err
}

func (o *Order) GetOrdersInProgress(ctx context.Context) ([]entity.Order, error) {
	return o.r.GetOrdersInProgress(ctx)

}

func (o *Order) UpdateOrder(ctx context.Context, dto UpdateOrderDTO) error {
	return o.r.UpdateOrder(ctx, repository.UpdateOrderDTO{
		Number:  dto.Number,
		Status:  dto.Status,
		Accrual: dto.Accrual,
		OrderID: dto.OrderID,
	})
}
