package handler

import (
	"github.com/DikosAs/GoAuthApi.git/src/storage"
	"github.com/gin-gonic/gin"
)

type Auth interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	CheckAuth(c *gin.Context)
}

type Handler struct {
	Auth
}

func NewHandler(storage *storage.Storage) *Handler {
	return &Handler{
		Auth: NewAuthHandler(storage),
	}
}
