package service

import (
	"github.com/DikosAs/GoAuthApi.git/src/repository"
)

type Authorisation interface {
	CreateUser(user repository.User) (repository.User, error)
	GenerateToken(username, password string) (*AuthTokens, error)
}

type Service struct {
	Authorisation
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorisation: NewAuthService(repo),
	}
}
