package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AntonioKichaev/internal/model"
	"github.com/AntonioKichaev/internal/service"
)

type AccountService interface {
	GetBalanceByID(ctx context.Context, dto *service.GetBalanceByIDDTO) (*model.UserBalance, error)
	Withdrawn(ctx context.Context, dto *service.WithdrawByUserIDDTO) error
	GetWithdraws(ctx context.Context, dto *service.WithdrawsByUserIDDTO) ([]model.Withdrawal, error)
}

type ResponseGetBalanceByUser struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type RequestWithdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type AccountHandle struct {
	ac AccountService
}

func NewAccountHandle(ac AccountService) *AccountHandle {
	return &AccountHandle{
		ac: ac,
	}
}

func (ah *AccountHandle) GetBalanceByUser(c *gin.Context) {

	userID, _ := c.Get("UserID")
	balance, err := ah.ac.GetBalanceByID(c.Request.Context(), &service.GetBalanceByIDDTO{
		UserID: userID.(uint),
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ResponseGetBalanceByUser{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	})
}

func (ah *AccountHandle) Withdraw(c *gin.Context) {
	req := RequestWithdraw{}
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	userID, _ := c.Get("UserID")

	err = ah.ac.Withdrawn(c.Request.Context(), &service.WithdrawByUserIDDTO{
		UserID:  userID.(uint),
		OrderID: req.Order,
		Amount:  req.Sum,
	})
	if err != nil {
		if errors.Is(err, service.ErrorNotFoundOrder) {
			//422 — неверный номер заказа;
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return

		}
		c.AbortWithStatus(http.StatusPaymentRequired)
		//c.AbortWithStatus(http.StatusInternalServerError)

		//402 — на счету недостаточно средств;
		//500 — внутренняя ошибка сервера.
		return
	}
	//200 — успешная обработка запроса;
	c.Status(http.StatusOK)

}

func (ah *AccountHandle) GetWithdraws(c *gin.Context) {

	//200 — успешная обработка запроса;
	//401 — пользователь не авторизован;
	//402 — на счету недостаточно средств;
	//422 — неверный номер заказа;
	//500 — внутренняя ошибка сервера.
	userID, _ := c.Get("UserID")
	wd, err := ah.ac.GetWithdraws(c.Request.Context(), &service.WithdrawsByUserIDDTO{
		UserID: userID.(uint),
	})
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, wd)

}
