package composite

import (
	"gorm.io/gorm"

	"github.com/AntonioKichaev/internal/controller"
	"github.com/AntonioKichaev/internal/repository"
	"github.com/AntonioKichaev/internal/service"
)

type Account struct {
	Repo    *repository.AccountRepository
	Service *service.Account
	Handle  *controller.AccountHandle
}

func NewAccount(db *gorm.DB, oa service.OrderAdapter, wd service.WithdrawnAdapter) *Account {
	ar := repository.NewAccountRepository(db)
	as := service.NewAccountService(ar, oa, wd)
	ah := controller.NewAccountHandle(as)

	return &Account{
		Repo:    ar,
		Service: as,
		Handle:  ah,
	}

}
