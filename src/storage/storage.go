package storage

import (
	"github.com/MarshMeg/auth-micro_service.git/src/types"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Auth interface {
	CreateUser(user types.User) (types.User, error)
	SetTokens(access *types.Token, refresh *types.Token) error
	GetUserIDByToken(t string) (string, int, error)
	UpdateUser(user types.User) error
	GetUsers(filters *types.User) ([]types.User, error)
	DeleteUser(user *types.User) error
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
