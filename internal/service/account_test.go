package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/AntonioKichaev/internal/entity"
	"github.com/AntonioKichaev/internal/mocks"
	"github.com/AntonioKichaev/internal/repository"
	"github.com/AntonioKichaev/internal/service"
)

func TestAccount_Withdrawn(t *testing.T) {
	repo := mocks.NewAccountRepository(t)
	oa := mocks.NewOrderAdapter(t)
	wd := mocks.NewWithdrawnAdapter(t)
	type args struct {
		ctx context.Context
		dto *service.WithdrawByUserIDDTO
	}
	type mockOn struct {
		method string
		args   []interface{}
	}

	tests := []struct {
		name                    string
		args                    args
		repoGetBalanceOn        mockOn
		repoUpdateAccOn         mockOn
		orderAdapterOn          mockOn
		wdAdapterOn             mockOn
		repoReturnArgs          []interface{}
		repoUpdateAccReturnArgs []interface{}
		orderAdapterReturnArgs  []interface{}
		wdReturnArgs            []interface{}
		wantErr                 bool
	}{

		{
			name: "not enought money",
			args: args{
				ctx: context.Background(),
				dto: &service.WithdrawByUserIDDTO{
					UserID:  1,
					OrderID: "1",
					Amount:  500,
				},
			},
			repoGetBalanceOn: mockOn{
				method: "GetBalanceByID",
				args: []interface{}{mock.Anything, repository.GetBalanceByIDDTO{
					UserID: 1,
				}},
			},
			repoReturnArgs: []interface{}{
				&entity.SaveBalance{
					UserID:  1,
					Current: 100.50,
				},
				nil,
			},
			wantErr: true,
		},
		{
			name: "enought money",
			args: args{
				ctx: context.Background(),
				dto: &service.WithdrawByUserIDDTO{
					UserID:  1,
					OrderID: "1",
					Amount:  50,
				},
			},
			repoGetBalanceOn: mockOn{
				method: "GetBalanceByID",
				args: []interface{}{mock.Anything, repository.GetBalanceByIDDTO{
					UserID: 1,
				}},
			},
			repoReturnArgs: []interface{}{
				&entity.SaveBalance{
					UserID:  1,
					Current: 100.50,
				},
				nil,
			},
			wantErr: false,
			orderAdapterOn: mockOn{
				method: "UploadOrderID",
				args: []interface{}{mock.Anything, &service.UploadOrderIDDTO{
					UserID: 1,
					Number: "1",
				}},
			},
			orderAdapterReturnArgs: []interface{}{
				&entity.Order{
					Number: "1",
					Model: gorm.Model{
						ID: 500,
					},
				},
				nil,
			},
			repoUpdateAccOn: mockOn{
				method: "UpdateAccount",
				args: []interface{}{mock.Anything, entity.SaveBalance{
					Current: 50.50,
					UserID:  1,
				},
				},
			},
			repoUpdateAccReturnArgs: []interface{}{
				nil,
			},
			wdAdapterOn: mockOn{
				method: "Create",
				args: []interface{}{mock.Anything, service.WithdrawnCreateDTO{
					OrderID: 500,
					UserID:  1,
					Sum:     50,
				},
				},
			},
			wdReturnArgs: []interface{}{
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.repoGetBalanceOn.args) > 0 || len(tt.repoReturnArgs) > 0 {
				repo.On(tt.repoGetBalanceOn.method, tt.repoGetBalanceOn.args...).Return(tt.repoReturnArgs...)
			}

			if len(tt.repoUpdateAccOn.args) > 0 || len(tt.repoUpdateAccReturnArgs) > 0 {
				repo.On(tt.repoUpdateAccOn.method, tt.repoUpdateAccOn.args...).Return(tt.repoUpdateAccReturnArgs...)
			}

			if len(tt.orderAdapterOn.args) > 0 || len(tt.orderAdapterReturnArgs) > 0 {
				oa.On(tt.orderAdapterOn.method, tt.orderAdapterOn.args...).Return(tt.orderAdapterReturnArgs...)

			}

			if len(tt.wdAdapterOn.args) > 0 || len(tt.wdReturnArgs) > 0 {
				wd.On(tt.wdAdapterOn.method, tt.wdAdapterOn.args...).Return(tt.wdReturnArgs...)
			}

			a := service.NewAccountService(repo, oa, wd)

			if err := a.Withdrawn(tt.args.ctx, tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("Withdrawn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
