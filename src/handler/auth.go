package handler

import (
	"fmt"
	"github.com/DikosAs/GoAuthApi.git/src/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Register(c *gin.Context) {
	var input repository.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	user, err := h.services.Authorisation.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tokens, err := h.services.Authorisation.GenerateToken(user.Username, user.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("X-Access-Token", tokens.AccessToken.Token, tokens.AccessToken.TTL, "/", "localhost", false, true)
	c.SetCookie("X-Refresh-Token", tokens.RefreshToken.Token, tokens.RefreshToken.TTL, "/", "localhost", false, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"tokens": tokens,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var input AuthLogin

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.services.Authorisation.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("X-Access-Token", tokens.AccessToken.Token, tokens.AccessToken.TTL, "/", "localhost", false, true)
	c.SetCookie("X-Refresh-Token", tokens.RefreshToken.Token, tokens.RefreshToken.TTL, "/", "localhost", false, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"tokens": tokens,
	})
}

func (h *Handler) CheckAuth(c *gin.Context) {
	fmt.Println(c.Params)
}
