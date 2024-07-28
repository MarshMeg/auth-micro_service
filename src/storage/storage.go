package storage

import (
	"github.com/DikosAs/GoAuthApi.git/src/storage/controllers"
	"github.com/DikosAs/GoAuthApi.git/src/types"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Auth interface {
	CreateUser(user types.User) (types.User, error)
	GetUser(u *types.User) (types.User, error)
	SetTokens(access *controllers.AccessToken, refresh *controllers.RefreshToken) error
	GetTokens(t string) (string, error)
}

type Cache interface {
}

type Storage struct {
	Auth
	Cache
}

func NewRepository(fileDB *gorm.DB, cache *redis.Client) *Storage {
	return &Storage{
		Auth:  NewAuthStorage(fileDB, cache),
		Cache: NewCacheDB(cache),
	}
}
