package storage

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Storage struct {
	Auth *AuthStorage
	User *UserStorage
}

func NewRepository(fileDB *gorm.DB, cache *redis.Client) *Storage {
	return &Storage{
		Auth: NewAuthStorage(fileDB, cache),
		User: NewUserStorage(fileDB, cache),
	}
}
