package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/DikosAs/auth-micro_service.git/src/types"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strings"
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

func (d *AuthStorage) GetUser(name string, id int) (types.User, error) {
	var user types.User
	if name != "" {
		err := d.db.Where("`username`=?", name).First(&user)
		return user, err.Error
	} else if id != 0 {
		err := d.db.Where("`id`=?", id).First(&user)
		return user, err.Error
	}
	return user, errors.New("user not found")
}

func (d *AuthStorage) SetTokens(access, refresh *types.Token) error {
	err := d.cache.Set(context.Background(), access.Token, fmt.Sprintf("access@%d", access.UserId), time.Duration(access.TTL)).Err()
	if err != nil {
		return err
	}
	err = d.cache.Set(context.Background(), refresh.Token, fmt.Sprintf("refresh@%d", refresh.UserId), time.Duration(refresh.TTL)).Err()
	return err
}

func (d *AuthStorage) GetUserIDByToken(t string) (string, int, error) {
	res, err := d.cache.Get(context.Background(), t).Result()
	if err != nil {
		return "", 0, err
	}
	parts := strings.Split(res, "@")
	if len(parts) != 2 {
		return "", 0, errors.New("incorrect key")
	}
	id, err := types.StrToInt(parts[1])
	return parts[0], id, err
}
