package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AntonioKichaev/config"
	"github.com/AntonioKichaev/internal/composite"
	"github.com/AntonioKichaev/internal/entity"
	"github.com/AntonioKichaev/internal/service"
	"github.com/AntonioKichaev/pkg/auth"
)

func Run() {

	cfg := config.New()
	route := gin.New()
	route.Use(gin.Logger())
	db, _ := gorm.Open(postgres.Open(cfg.GetDatabaseURI()), &gorm.Config{})

	authMid := auth.NewAuthMid("322")
	// Migrate the schema
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Order{},
		&entity.SaveBalance{},
		&entity.Withdrawn{},
	)
	if err != nil {
		logrus.Fatal(err)
	}

	wc := composite.NewWithdrawn(db)
	oc := composite.NewOrder(db)
	ah := composite.NewAccount(db, oc.Service, wc.Service)
	uc := composite.NewAuth(db, authMid, ah.Service)
	gr := route.Group("/api")
	{

		gr.POST("user/register", uc.Handle.Register)
		gr.POST("user/login", uc.Handle.Login)

	}

	{
		//auth required

		authGr := gr.Group("/user")
		authGr.Use(auth.AuthMiddleware(authMid))

		authGr.POST("/orders", oc.Handle.UploadOrderID)
		authGr.GET("/orders", oc.Handle.GetOrderByUser)
		//

		authGr.GET("/balance", ah.Handle.GetBalanceByUser)
		authGr.POST("/balance/withdraw", ah.Handle.Withdraw)

		authGr.GET("/withdrawals", ah.Handle.GetWithdraws)
	}

	ac := service.NewAccrual(oc.Service, ah.Service, cfg.GetAccrualSystemAddress())
	go ac.Start(context.Background())
	logrus.Info("start listen ", cfg.GetRunAddress())
	err = http.ListenAndServe(cfg.GetRunAddress(), route)
	logrus.Error(err)
}
