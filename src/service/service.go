package service

import (
	"github.com/DikosAs/GoAuthApi.git/src/repository"
	"github.com/DikosAs/GoAuthApi.git/src/repository/objects"
)

type Authorisation interface {
	CreateUser(user objects.User) (int, error)
	GenerateToken(username, password string) (string, error)
}

type Service struct {
	Authorisation
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorisation: NewAuthService(repo),
	}
}
