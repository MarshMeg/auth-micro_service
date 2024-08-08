package storage

import (
	"context"
	"fmt"
	"github.com/MarshMeg/auth-micro_service.git/src/types/auth"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

type AuthStorage struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewAuthStorage(db *gorm.DB, redis *redis.Client) *AuthStorage {
	return &AuthStorage{
		db:    db,
		cache: redis,
	}
}

func (d *AuthStorage) SetTokens(access, refresh *auth.Token) error {
	err := d.cache.Set(context.Background(), access.Token, fmt.Sprintf("access@%d", access.UserId), time.Duration(access.TTL)).Err()
	if err != nil {
		return err
	}
	err = d.cache.Set(context.Background(), refresh.Token, fmt.Sprintf("refresh@%d", refresh.UserId), time.Duration(refresh.TTL)).Err()
	return err
}
