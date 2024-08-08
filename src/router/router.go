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

func (r *Router) InitRoutes(mode, prefix string) *gin.Engine {
	gin.SetMode(mode)
	router := gin.New()

	root := router.Group(prefix)
	{
		root.GET("/status", r.CheckStatus)
		auth := root.Group("/auth")
		{
			auth.POST("/register", r.handler.Auth.Register)
			auth.POST("/login", r.handler.Auth.Login)
			auth.GET("/check_auth", r.handler.Auth.CheckAuth)
		}
		users := root.Group("/users")
		{
			users.GET("/:id", r.handler.User.GetUser)
			users.PATCH("/:id", r.handler.User.PatchUser)
			users.DELETE("/:id", r.handler.User.DeleteUser)
			users.GET("/", r.handler.User.GetUsers)
		}
		groups := root.Group("/groups")
		{
			groups.GET("/", r.CheckStatus)
			groups.POST("/", r.CheckStatus)
			groups.PATCH("/", r.CheckStatus)
			groups.DELETE("/", r.CheckStatus)
		}
	}

	return router
}

func (r *Router) CheckStatus(c *gin.Context) {
	c.JSON(200, map[string]interface{}{})
}
