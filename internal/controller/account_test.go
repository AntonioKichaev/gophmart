package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/AntonioKichaev/internal/controller/mocks"
	"github.com/AntonioKichaev/internal/model"
	"github.com/AntonioKichaev/internal/service"
)

func TestAccountHandle_GetBalanceByUser(t *testing.T) {
	m := mocks.NewAccountService(t)
	var UserID uint = 1
	expResp := ResponseGetBalanceByUser{
		Current:   100,
		Withdrawn: 155.1,
	}
	m.On("GetBalanceByID", mock.Anything, &service.GetBalanceByIDDTO{
		UserID: UserID,
	}).Return(&model.UserBalance{
		Current:   expResp.Current,
		Withdrawn: expResp.Withdrawn,
	}, nil)

	h := NewAccountHandle(m)
	gin.SetMode(gin.TestMode)
	route := gin.New()
	route.GET("/balance", func(context *gin.Context) {
		context.Set("UserID", UserID)
		h.GetBalanceByUser(context)
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/balance", nil)
	route.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", w.Code)
	}
	resp := ResponseGetBalanceByUser{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, expResp, resp)
	m.AssertExpectations(t)

}
