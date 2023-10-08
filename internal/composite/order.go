package composite

import (
	"gorm.io/gorm"

	"github.com/AntonioKichaev/internal/controller"
	"github.com/AntonioKichaev/internal/repository"
	"github.com/AntonioKichaev/internal/service"
)

type Order struct {
	Repo    *repository.OrderRepository
	Service *service.Order
	Handle  *controller.OrderHandle
}

func NewOrder(db *gorm.DB) *Order {
	orderRepo := repository.NewOrder(db)
	orderService := service.NewOrder(orderRepo)
	oh := controller.NewOrderHandle(orderService)
	return &Order{
		Repo:    orderRepo,
		Service: orderService,
		Handle:  oh,
	}

}
