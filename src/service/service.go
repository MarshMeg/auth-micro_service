package service

import (
	"go_back/src/repository"
	"go_back/src/repository/objects"
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
