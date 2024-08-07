package router

import (
	"github.com/MarshMeg/auth-micro_service.git/src/handler"
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
			user := v1.Group("/user")
			{
				user.POST("/register", r.handler.Register)
				user.POST("/login", r.handler.Login)
				user.GET("/check_auth", r.handler.CheckAuth)
				user.GET("/users", r.handler.GetUsers)

				user.GET("/:id", r.handler.GetUserByID)
				user.PATCH("/", r.handler.PatchUser)
				user.DELETE("/:id", r.handler.DeleteUser)
			}

		}
	}

	return router
}
