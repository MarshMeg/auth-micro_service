package service

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/DikosAs/GoAuthApi.git/src/repository"
	"io"
	"time"
)

const (
	accessTokenTTL  = 12 * time.Hour
	refreshTokenTTL = 30 * (24 * time.Hour)
)

type AuthService struct {
	repo repository.Authorisation
}

type Token struct {
	token string
	ttl   int
}

type AuthTokens struct {
	AccessToken  *repository.AccessToken
	RefreshToken *repository.RefreshToken
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user repository.User) (repository.User, error) {
	user.Password = passwdHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (*AuthTokens, error) {
	user, _ := s.repo.GetUser(username)

	accessToken := make([]byte, 128)
	refreshToken := make([]byte, 128)

	if _, err := io.ReadFull(rand.Reader, accessToken); err != nil {
		return &AuthTokens{}, err
	}
	if _, err := io.ReadFull(rand.Reader, refreshToken); err != nil {
		return &AuthTokens{}, err
	}

	tokens := &AuthTokens{
		AccessToken: &repository.AccessToken{
			UserID: user.Id,
			Token:  base64.RawURLEncoding.EncodeToString(accessToken),
			TTL:    int(accessTokenTTL),
		},
		RefreshToken: &repository.RefreshToken{
			UserID: user.Id,
			Token:  base64.RawURLEncoding.EncodeToString(refreshToken),
			TTL:    int(refreshTokenTTL),
		},
	}

	if err := s.repo.SetTokens(tokens.AccessToken, tokens.RefreshToken); err != nil {
		return &AuthTokens{}, err
	}
	return tokens, nil
}

func passwdHash(passwd string) string {
	var hash = sha512.New()
	hash.Write([]byte(passwd))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
