package service

import (
	"crypto/sha512"
	"fmt"
	"github.com/DikosAs/GoAuthApi.git/src/repository"
	"github.com/DikosAs/GoAuthApi.git/src/repository/objects"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	salt       = "hjqrhjqw124617django"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type AuthService struct {
	repo repository.Authorisation
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user objects.User) (int, error) {
	user.Password = passwdHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, passwdHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func passwdHash(passwd string) string {
	var hash = sha512.New()
	hash.Write([]byte(passwd))

	return fmt.Sprintf("%x", string(hash.Sum([]byte(salt))))
}
