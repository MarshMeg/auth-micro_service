package router

import (
	"github.com/DikosAs/auth-micro_service.git/src/handler"
	"github.com/gin-gonic/gin"
)

type Router struct {
	handler *handler.Handler
}

func NewRouter(handler *handler.Handler) *Router {
	return &Router{handler: handler}
}

func (r *Router) InitRoutes(mode string) *gin.Engine {
	gin.SetMode(mode)
	router := gin.New()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/register", r.handler.Register)
				auth.POST("/login", r.handler.Login)
				auth.GET("/check_auth", r.handler.CheckAuth)
				auth.GET("/users", r.handler.GetUsers)

				auth.DELETE("/user/:id", r.handler.DeleteUser)
			}
		}
	}

	return router
}
