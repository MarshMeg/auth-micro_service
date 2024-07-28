package handler

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/DikosAs/auth-micro_service.git/src/storage"
	"github.com/DikosAs/auth-micro_service.git/src/types"
	"github.com/DikosAs/auth-micro_service.git/src/types/request"
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
	AccessToken  *types.Token
	RefreshToken *types.Token
}

const (
	accessTokenTTL  = 12 * time.Hour
	refreshTokenTTL = 30 * (24 * time.Hour)
)

func (h *AuthHandler) Register(c *gin.Context) {
	var input types.User

	if err := c.BindJSON(&input); err != nil || input.Username == "" || input.Password == "" {
		newErrorResponse(c, http.StatusBadRequest, "Invalid body")
		return
	}

	if _, err := h.storage.Auth.GetUser(input.Username, 0); err == nil {
		newErrorResponse(c, http.StatusBadRequest, "User is already registered")
		return
	}

	input.Password = passwdHash(input.Password)
	user, err := h.storage.Auth.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tokens, err := h.generateTokens(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("X-Access-Token", tokens.AccessToken.Token, tokens.AccessToken.TTL, "/", "localhost", false, true)
	c.SetCookie("X-Refresh-Token", tokens.RefreshToken.Token, tokens.RefreshToken.TTL, "/", "localhost", false, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"username": user.Username,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input AuthLogin

	if err := c.BindJSON(&input); err != nil || input.Username == "" || input.Password == "" {
		newErrorResponse(c, http.StatusBadRequest, "Invalid body")
		return
	}

	user, err := h.storage.Auth.GetUser(input.Username, 0)

	tokens, err := h.generateTokens(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("X-Access-Token", tokens.AccessToken.Token, tokens.AccessToken.TTL, "/", "localhost", false, true)
	c.SetCookie("X-Refresh-Token", tokens.RefreshToken.Token, tokens.RefreshToken.TTL, "/", "localhost", false, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"username": user.Username,
	})
}

func (h *AuthHandler) CheckAuth(c *gin.Context) {
	tokenType, userId, _ := h.getAuth(c)

	user, _ := h.storage.Auth.GetUser("", userId)
	user.Password = "<secret>"

	c.JSON(http.StatusOK, map[string]interface{}{
		"user":       user,
		"token_type": tokenType,
	})
}

func (h *AuthHandler) PatchUser(c *gin.Context) {
	_, userId, initializer := h.getAuth(c)
	if initializer == request.Service {
		newErrorResponse(c, http.StatusForbidden, "The services do not have access to user accounts")
		return
	}

	var input types.User
	if err := c.BindJSON(&input); err != nil || input.Username == "" {
		newErrorResponse(c, http.StatusBadRequest, "Invalid body")
		return
	}

	if userId != input.Id {
		newErrorResponse(c, http.StatusForbidden, "You do not have access to this account")
	}

	input.Password = passwdHash(input.Password)
	err := h.storage.UpdateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})
}

func passwdHash(passwd string) string {
	var hash = sha512.New()
	hash.Write([]byte(passwd))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (h *AuthHandler) generateTokens(user types.User) (*AuthTokens, error) {
	accessToken := make([]byte, 128)
	refreshToken := make([]byte, 128)

	if _, err := io.ReadFull(rand.Reader, accessToken); err != nil {
		return &AuthTokens{}, err
	}
	if _, err := io.ReadFull(rand.Reader, refreshToken); err != nil {
		return &AuthTokens{}, err
	}

	tokens := &AuthTokens{
		AccessToken: &types.Token{
			UserId: user.Id,
			Token:  base64.RawURLEncoding.EncodeToString(accessToken),
			TTL:    int(accessTokenTTL),
		},
		RefreshToken: &types.Token{
			UserId: user.Id,
			Token:  base64.RawURLEncoding.EncodeToString(refreshToken),
			TTL:    int(refreshTokenTTL),
		},
	}

	if err := h.storage.Auth.SetTokens(tokens.AccessToken, tokens.RefreshToken); err != nil {
		return &AuthTokens{}, err
	}
	return tokens, nil
}

func (h *AuthHandler) getAuth(c *gin.Context) (string, int, int) {
	var token string
	var initializer int
	switch c.GetHeader("X-Real-IP") {
	case "service":
		token = c.Query("token")
		initializer = request.Service
		if token == "" {
			newErrorResponse(c, http.StatusBadRequest, "Invalid params")
			return "", 0, initializer
		}
	default:
		token, _ = c.Cookie("X-Access-Token")
		initializer = request.User
		if token == "" {
			newErrorResponse(c, http.StatusUnauthorized, "Token not found in \"X-Access-Token\" header. You not authenticated")
			return "", 0, initializer
		}
	}
	tokenType, userId, err := h.storage.Auth.GetUserIDByToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Not authenticated")
		return "", 0, initializer
	}

	return tokenType, userId, initializer
}
