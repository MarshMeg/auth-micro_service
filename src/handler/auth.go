package handler

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/MarshMeg/auth-micro_service.git/src/storage"
	"github.com/MarshMeg/auth-micro_service.git/src/types/auth"
	"github.com/MarshMeg/auth-micro_service.git/src/types/user"
	"github.com/MarshMeg/auth-micro_service.git/src/types/user/roles"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

type AuthHandler struct {
	storage *storage.Storage
}

func NewAuthHandler(storage *storage.Storage) *AuthHandler {
	return &AuthHandler{storage: storage}
}

type AuthLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthTokens struct {
	AccessToken  *auth.Token
	RefreshToken *auth.Token
}

const (
	accessTokenTTL  = 12 * time.Hour
	refreshTokenTTL = 30 * (24 * time.Hour)
)

func (h *AuthHandler) Register(c *gin.Context) {
	var input user.User

	if err := c.BindJSON(&input); err != nil || input.Username == "" || input.Password == "" {
		newErrorResponse(c, http.StatusBadRequest, "Invalid body")
		return
	}
	input.RoleName = roles.Member().RoleName

	if _, err := h.storage.User.GetUsers(&user.User{Username: input.Username}); err == nil {
		newErrorResponse(c, http.StatusBadRequest, "User is already registered")
		return
	}

	input.Password = passwdHash(input.Password)
	createdUser, err := h.storage.User.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tokens, err := generateTokens(h.storage.Auth, createdUser)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("X-Access-Token", tokens.AccessToken.Token, tokens.AccessToken.TTL, "/", "localhost", false, true)
	c.SetCookie("X-Refresh-Token", tokens.RefreshToken.Token, tokens.RefreshToken.TTL, "/", "localhost", false, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"username": createdUser.Username,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input AuthLogin
	if err := c.BindJSON(&input); err != nil || input.Username == "" || input.Password == "" {
		newErrorResponse(c, http.StatusBadRequest, "Invalid body")
		return
	}

	getUser, err := h.storage.User.GetUser(&user.User{Username: input.Username})

	tokens, err := generateTokens(h.storage.Auth, getUser)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("X-Access-Token", tokens.AccessToken.Token, tokens.AccessToken.TTL, "/", "localhost", false, true)
	c.SetCookie("X-Refresh-Token", tokens.RefreshToken.Token, tokens.RefreshToken.TTL, "/", "localhost", false, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"username": getUser.Username,
	})
}

func (h *AuthHandler) CheckAuth(c *gin.Context) {
	_, authUser := NewAuthService(h.storage).CheckAuth(c, 0, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"user": authUser.Return(),
	})
}

func passwdHash(passwd string) string {
	var hash = sha512.New()
	hash.Write([]byte(passwd))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func generateTokens(storage *storage.AuthStorage, user user.User) (*AuthTokens, error) {
	accessToken := make([]byte, 128)
	refreshToken := make([]byte, 128)

	if _, err := io.ReadFull(rand.Reader, accessToken); err != nil {
		return &AuthTokens{}, err
	}
	if _, err := io.ReadFull(rand.Reader, refreshToken); err != nil {
		return &AuthTokens{}, err
	}

	tokens := &AuthTokens{
		AccessToken: &auth.Token{
			UserId: user.Id,
			Token:  base64.RawURLEncoding.EncodeToString(accessToken),
			TTL:    int(accessTokenTTL),
		},
		RefreshToken: &auth.Token{
			UserId: user.Id,
			Token:  base64.RawURLEncoding.EncodeToString(refreshToken),
			TTL:    int(refreshTokenTTL),
		},
	}

	if err := storage.SetTokens(tokens.AccessToken, tokens.RefreshToken); err != nil {
		return &AuthTokens{}, err
	}
	return tokens, nil
}
