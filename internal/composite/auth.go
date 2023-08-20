package composite

import (
	"gorm.io/gorm"

	"github.com/AntonioKichaev/internal/controller"
	"github.com/AntonioKichaev/internal/repository"
	"github.com/AntonioKichaev/internal/service"
	"github.com/AntonioKichaev/pkg/auth"
)

type Auth struct {
	Repo    *repository.AuthRepository
	Service *service.AuthService
	Handle  *controller.Auth
}

func NewAuth(db *gorm.DB, mid *auth.AuthMid, as service.AccountServiceAdapter) *Auth {
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, as)
	AuthHandle := controller.NewAuthHandle(authService, mid)
	return &Auth{
		Repo:    authRepo,
		Service: authService,
		Handle:  AuthHandle,
	}

}
