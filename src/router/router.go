package router

import (
	"github.com/DikosAs/GoAuthApi.git/src/handler"
	"github.com/gin-gonic/gin"
)

type Router struct {
	handler *handler.Handler
}

func NewRouter(handler *handler.Handler) *Router {
	return &Router{handler: handler}
}

func (r *Router) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/register", r.handler.Register)
				auth.POST("/login", r.handler.Login)
				//auth.GET("/check_auth", r.handler.CheckAuth)
			}
		}
	}

	return router
}
