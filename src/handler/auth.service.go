package handler

import (
	"github.com/MarshMeg/auth-micro_service.git/src/storage"
	"github.com/MarshMeg/auth-micro_service.git/src/types/user"
	"github.com/MarshMeg/auth-micro_service.git/src/types/user/roles"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthService struct {
	storage *storage.Storage
}

func NewAuthService(storage *storage.Storage) *AuthService {
	return &AuthService{storage: storage}
}

func (s *AuthService) CheckAuth(c *gin.Context, minRoleLvl int, accessCritical bool) (bool, user.User) {
	var authUser user.User

	token, _ := c.Cookie("X-Access-Token")
	if token == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Token not found in \"X-Access-Token\" header. You not authenticated")
		return false, authUser
	}

	_, userId, err := s.storage.User.GetUserIDByToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Not authenticated")
		return false, authUser
	}
	authUser, err = s.storage.User.GetUser(&user.User{Id: userId})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "User not found by token")
		return false, authUser
	}

	list := roles.List()
	if !(list.GetRole(&roles.UserRole{RoleName: authUser.RoleName}).AccessLvl > minRoleLvl) {
		if accessCritical {
			newErrorResponse(c, http.StatusForbidden, "You do not have access")
		}
		return false, authUser
	}
	return true, authUser
}
