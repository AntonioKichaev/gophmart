package service

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/AntonioKichaev/internal/entity"
)

type OrderCheckAdapter interface {
	GetOrdersInProgress(ctx context.Context) ([]entity.Order, error)
	UpdateOrder(ctx context.Context, dto UpdateOrderDTO) error
}
type AccountAdapter interface {
	AddPoint(ctx context.Context, dto AddPointDTO) error
}
type Accrual struct {
	oa                 OrderCheckAdapter
	aa                 AccountAdapter
	client             *resty.Client
	accrualExternalURL string
}

func NewAccrual(oa OrderCheckAdapter, aa AccountAdapter, accrualAddr string) *Accrual {
	return &Accrual{
		oa:                 oa,
		aa:                 aa,
		client:             resty.New(),
		accrualExternalURL: accrualAddr,
	}
}

type ResponseAccrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (a *Accrual) CheckOrder(ctx context.Context, order entity.Order) {
	result, err := url.JoinPath(a.accrualExternalURL, "/api/orders/", order.Number)
	if err != nil {
		logrus.Error("a.CheckOrder()", err)
		return
	}
	resp, err := a.client.R().Get(result)
	if err != nil {
		return
	}
	ra := ResponseAccrual{}
	err = json.Unmarshal(resp.Body(), &ra)
	if err != nil {
		logrus.Error("json.Unmarshal", err, string(resp.Body()))
		return
	}

	err = a.oa.UpdateOrder(ctx, UpdateOrderDTO{
		Number:  ra.Order,
		Status:  ra.Status,
		Accrual: ra.Accrual,
		OrderID: order.ID,
	})

	if err != nil {
		logrus.Error("a.oa.UpdateOrder()", err)
	}
	err = a.aa.AddPoint(ctx, AddPointDTO{
		Point:  ra.Accrual,
		UserID: order.UserID,
	})
	if err != nil {
		logrus.Error("a.aa.AddPoint()", err)
	}

}

func (a *Accrual) Start(ctx context.Context) {
	for ; ; time.Sleep(time.Second) {
		select {
		case <-ctx.Done():
			return
		default:
			orders, err := a.oa.GetOrdersInProgress(ctx)
			if err != nil {
				logrus.Error("error a.oa.GetOrdersInProgress()", err)
				continue
			}
			if len(orders) == 0 {
				continue
			}

			for _, order := range orders {
				a.CheckOrder(ctx, order)
			}

		}

	}
}
