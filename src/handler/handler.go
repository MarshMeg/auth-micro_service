package handler

import (
	"github.com/MarshMeg/auth-micro_service.git/src/storage"
)

type Handler struct {
	Auth *AuthHandler
	User *UserHandler
}

func NewHandler(storage *storage.Storage) *Handler {
	return &Handler{
		Auth: NewAuthHandler(storage),
		User: NewUserHandler(storage),
	}
}
