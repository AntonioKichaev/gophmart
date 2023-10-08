package composite

import (
	"gorm.io/gorm"

	"github.com/AntonioKichaev/internal/repository"
	"github.com/AntonioKichaev/internal/service"
)

type Withdrawn struct {
	Repo    *repository.WithdrawnRepository
	Service *service.WithdrawnService
}

func NewWithdrawn(db *gorm.DB) *Withdrawn {
	r := repository.NewWithdrawnRepository(db)
	s := service.NewWithdrawnService(r)
	return &Withdrawn{
		Repo:    r,
		Service: s,
	}
}
