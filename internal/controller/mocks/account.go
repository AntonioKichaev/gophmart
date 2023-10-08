// Code generated by mockery v2.30.16. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/AntonioKichaev/internal/model"

	service "github.com/AntonioKichaev/internal/service"
)

// AccountService is an autogenerated mock type for the AccountService type
type AccountService struct {
	mock.Mock
}

// GetBalanceByID provides a mock function with given fields: ctx, dto
func (_m *AccountService) GetBalanceByID(ctx context.Context, dto *service.GetBalanceByIDDTO) (*model.UserBalance, error) {
	ret := _m.Called(ctx, dto)

	var r0 *model.UserBalance
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *service.GetBalanceByIDDTO) (*model.UserBalance, error)); ok {
		return rf(ctx, dto)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *service.GetBalanceByIDDTO) *model.UserBalance); ok {
		r0 = rf(ctx, dto)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserBalance)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *service.GetBalanceByIDDTO) error); ok {
		r1 = rf(ctx, dto)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWithdraws provides a mock function with given fields: ctx, dto
func (_m *AccountService) GetWithdraws(ctx context.Context, dto *service.WithdrawsByUserIDDTO) ([]model.Withdrawal, error) {
	ret := _m.Called(ctx, dto)

	var r0 []model.Withdrawal
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *service.WithdrawsByUserIDDTO) ([]model.Withdrawal, error)); ok {
		return rf(ctx, dto)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *service.WithdrawsByUserIDDTO) []model.Withdrawal); ok {
		r0 = rf(ctx, dto)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Withdrawal)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *service.WithdrawsByUserIDDTO) error); ok {
		r1 = rf(ctx, dto)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Withdrawn provides a mock function with given fields: ctx, dto
func (_m *AccountService) Withdrawn(ctx context.Context, dto *service.WithdrawByUserIDDTO) error {
	ret := _m.Called(ctx, dto)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *service.WithdrawByUserIDDTO) error); ok {
		r0 = rf(ctx, dto)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAccountService creates a new instance of AccountService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccountService(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccountService {
	mock := &AccountService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
