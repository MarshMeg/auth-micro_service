package repository

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Authorisation interface {
	CreateUser(user User) (User, error)
	GetUser(username string) (User, error)
	SetTokens(access *AccessToken, refresh *RefreshToken) error
}

type Cache interface {
}

type Repository struct {
	Authorisation
	Cache
}

func NewRepository(fileDB *gorm.DB, cache *redis.Client) *Repository {
	return &Repository{
		Authorisation: NewAuthDB(fileDB, cache),
		Cache:         NewCacheDB(cache),
	}
}
