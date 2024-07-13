package handler

import (
	"github.com/DikosAs/GoAuthApi.git/src/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/register", h.Register)
				auth.POST("/login", h.Login)
			}

			//mail := v1.Group("/mail")
			//{
			//	mail.POST("/sand")
			//}
		}
	}

	return router
}
