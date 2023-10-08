package controller

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"

	"github.com/AntonioKichaev/internal/entity"
	"github.com/AntonioKichaev/internal/service"
)

type OrderService interface {
	UploadOrderID(ctx context.Context, dto *service.UploadOrderIDDTO) (*entity.Order, error)
	GetOrdersByUser(ctx context.Context, dto *service.GetOrdersByUserIDDTO) ([]entity.Order, error)
}

type OrderHandle struct {
	o OrderService
}

func NewOrderHandle(o OrderService) *OrderHandle {
	return &OrderHandle{o: o}
}

func (oh OrderHandle) UploadOrderID(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return

	}
	number := string(data)
	userID, _ := c.Get("UserID")

	if err := isValidLunaCheckNumber(number); err != nil {
		// 422 неверный формат заказа
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	_, err = oh.o.UploadOrderID(c.Request.Context(), &service.UploadOrderIDDTO{
		Number: number,
		UserID: userID.(uint),
	})
	if err != nil {
		if errors.Is(err, service.ErrOrderExists) {
			//200 — номер заказа уже был загружен этим пользователем;
			c.Status(http.StatusOK)
			return
		}

		// 409 номер заказа уже был загружен другим пользователем
		if errors.Is(err, service.ErrOrderUploadedOtherUser) {
			c.Status(http.StatusConflict)

			return
		}
		// 500
		c.Status(http.StatusInternalServerError)

		return
	}
	//202 — новый номер заказа принят в обработку;
	c.Status(http.StatusAccepted)
}

func (oh OrderHandle) GetOrderByUser(c *gin.Context) {

	userID, _ := c.Get("UserID")

	orders, err := oh.o.GetOrdersByUser(c.Request.Context(), &service.GetOrdersByUserIDDTO{
		UserID: userID.(uint),
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		// 500
		return
	}

	if len(orders) == 0 {
		//204 — нет данных для ответа
		c.Status(http.StatusNotFound)
		return
	}

	//200 — orders exits
	c.JSON(http.StatusOK, orders)
}

func isValidLunaCheckNumber(s string) error {

	return goluhn.Validate(s)
}
