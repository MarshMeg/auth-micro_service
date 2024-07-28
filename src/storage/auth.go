package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/DikosAs/GoAuthApi.git/src/storage/controllers"
	"github.com/DikosAs/GoAuthApi.git/src/types"
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

func (d *AuthStorage) CreateUser(user types.User) (types.User, error) {
	result := d.db.Create(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (d *AuthStorage) GetUser(input *types.User) (types.User, error) {
	var user types.User
	if input.Username != "" {
		err := d.db.Where("`username`=?", input.Username).First(&user)
		return user, err.Error
	}
	if input.Id != 0 {
		err := d.db.Where("`userId`=?", input.Id).First(&user)
		return user, err.Error
	}
	return user, errors.New("user not found")
}

func (d *AuthStorage) SetTokens(access *controllers.AccessToken, refresh *controllers.RefreshToken) error {
	err := d.cache.Set(context.Background(), access.Token, fmt.Sprintf("access_key.%d", access.UserID), time.Duration(access.TTL)).Err()
	if err != nil {
		return err
	}
	err = d.cache.Set(context.Background(), refresh.Token, fmt.Sprintf("refresh_key.%d", refresh.UserID), time.Duration(refresh.TTL)).Err()
	return err
}

func (d *AuthStorage) GetTokens(t string) (string, error) {
	res, err := d.cache.Get(context.Background(), fmt.Sprintf("%s", t)).Result()
	return res, err
}
