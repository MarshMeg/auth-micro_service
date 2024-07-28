package handler

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/DikosAs/GoAuthApi.git/src/storage"
	"github.com/DikosAs/GoAuthApi.git/src/storage/controllers"
	"github.com/DikosAs/GoAuthApi.git/src/types"
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

type Token struct {
	token string
	ttl   int
}

type AuthTokens struct {
	AccessToken  *controllers.AccessToken
	RefreshToken *controllers.RefreshToken
}

const (
	accessTokenTTL  = 12 * time.Hour
	refreshTokenTTL = 30 * (24 * time.Hour)
)

func (h *AuthHandler) Register(c *gin.Context) {
	var input types.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
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

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	user, err := h.storage.Auth.GetUser(&types.User{
		Username: input.Username,
		Password: passwdHash(input.Password),
	})

	tokens, err := h.generateTokens(user)
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

//func (h *AuthHandler) CheckAuth(c *gin.Context) {
//	token := c.Query("token")
//	if token == "" {
//		newErrorResponse(c, http.StatusBadRequest, "invalid token")
//		return
//	}
//
//	user, err := h.services.VerifyToken(token)
//	if err != nil {
//		newErrorResponse(c, http.StatusInternalServerError, "token not found in database")
//		return
//	}
//
//	c.JSON(http.StatusOK, map[string]interface{}{
//		"user": user,
//	})
//}

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
		AccessToken: &controllers.AccessToken{
			UserID: user.Id,
			Token:  base64.RawURLEncoding.EncodeToString(accessToken),
			TTL:    int(accessTokenTTL),
		},
		RefreshToken: &controllers.RefreshToken{
			UserID: user.Id,
			Token:  base64.RawURLEncoding.EncodeToString(refreshToken),
			TTL:    int(refreshTokenTTL),
		},
	}

	if err := h.storage.Auth.SetTokens(tokens.AccessToken, tokens.RefreshToken); err != nil {
		return &AuthTokens{}, err
	}
	return tokens, nil
}
