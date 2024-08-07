package handler

import (
	"github.com/MarshMeg/auth-micro_service.git/src/storage"
	"github.com/gin-gonic/gin"
)

type Auth interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	CheckAuth(c *gin.Context)
	PatchUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUserByID(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type Handler struct {
	Auth
}

func NewHandler(storage *storage.Storage) *Handler {
	return &Handler{
		Auth: NewAuthHandler(storage),
	}
}
