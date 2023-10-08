package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/AntonioKichaev/internal/entity"
)

type CreateOrderDTO struct {
	UserID     uint
	Number     string
	Status     string
	UploadedAt time.Time
}

type GetOrdersDTO struct {
	UserID uint
}

type GetOrderByNumberDTO struct {
	Number string
}
type UpdateOrderDTO struct {
	Number  string
	Status  string
	Accrual float64
	OrderID uint
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrder(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (o *OrderRepository) Create(ctx context.Context, dto CreateOrderDTO) (entity.Order, error) {
	order := &entity.Order{
		UserID:     dto.UserID,
		Number:     dto.Number,
		Status:     dto.Status,
		UploadedAt: dto.UploadedAt,
	}
	tx := o.db.Model(entity.Order{}).Create(order)
	return *order, tx.Error
}

func (o *OrderRepository) GetOrders(ctx context.Context, dto GetOrdersDTO) []entity.Order {
	var orders []entity.Order
	o.db.Model(entity.Order{}).Where("user_id = ?", dto.UserID).Find(&orders)
	return orders
}

func (o *OrderRepository) GetOrder(ctx context.Context, dto GetOrderByNumberDTO) (*entity.Order, error) {
	var order entity.Order
	tx := o.db.Model(entity.Order{}).Where("number = ?", dto.Number).First(&order)
	return &order, tx.Error
}

func (o *OrderRepository) UpdateOrder(ctx context.Context, dto UpdateOrderDTO) error {
	order := entity.Order{
		Number:  dto.Number,
		Status:  dto.Status,
		Accrual: dto.Accrual,
	}
	tx := o.db.Model(entity.Order{}).Where("id = ?", dto.OrderID).Updates(&order)
	return tx.Error
}
func (o *OrderRepository) GetOrdersInProgress(ctx context.Context) ([]entity.Order, error) {
	var orders []entity.Order
	tx := o.db.Model(entity.Order{}).Where("status IN ?", []string{"REGISTERED", "PROCESSING", "NEW"}).Find(&orders)

	return orders, tx.Error
}
