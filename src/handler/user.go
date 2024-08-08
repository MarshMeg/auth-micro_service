package handler

import (
	"github.com/MarshMeg/auth-micro_service.git/src/storage"
	"github.com/MarshMeg/auth-micro_service.git/src/types"
	"github.com/MarshMeg/auth-micro_service.git/src/types/user"
	"github.com/MarshMeg/auth-micro_service.git/src/types/user/roles"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	storage *storage.Storage
}

func NewUserHandler(storage *storage.Storage) *UserHandler {
	return &UserHandler{storage: storage}
}

func (h *UserHandler) PatchUser(c *gin.Context) {
	accessBool, authUser := NewAuthService(h.storage).CheckAuth(c, roles.List().Service.AccessLvl, false)

	var input user.User
	if err := c.BindJSON(&input); err != nil || input.Username == "" {
		newErrorResponse(c, http.StatusBadRequest, "Invalid body")
		return
	}

	if !accessBool && authUser.Id != input.Id {
		newErrorResponse(c, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	input.Password = passwdHash(input.Password)
	err := h.storage.User.UpdateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error")
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	NewAuthService(h.storage).CheckAuth(c, roles.List().Service.AccessLvl, true)

	users, err := h.storage.User.GetUsers(&user.User{
		Id:       types.StrToInt(c.Query("id")),
		Username: c.Query("username"),
		RoleName: c.Query("role"),
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	NewAuthService(h.storage).CheckAuth(c, roles.List().Service.AccessLvl, true)

	val, _ := c.Params.Get("id")
	getUser, err := h.storage.User.GetUser(&user.User{Id: types.StrToInt(val)})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getUser)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	NewAuthService(h.storage).CheckAuth(c, roles.List().Admin.AccessLvl, true)

	val, _ := c.Params.Get("id")
	getUser, err := h.storage.User.GetUser(&user.User{Id: types.StrToInt(val)})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.storage.User.DeleteUser(&getUser)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}
