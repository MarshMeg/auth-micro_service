package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

type AuthDB struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewAuthDB(db *gorm.DB, rediska *redis.Client) *AuthDB {
	return &AuthDB{db: db, cache: rediska}
}

func (d *AuthDB) CreateUser(user User) (User, error) {
	result := d.db.Create(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (d *AuthDB) GetUser(username string) (User, error) {
	var user User
	err := d.db.Where("`username`=?", username).First(&user)
	return user, err.Error
}

func (d *AuthDB) SetTokens(access *AccessToken, refresh *RefreshToken) error {
	err := d.cache.Set(context.Background(), access.Token, fmt.Sprintf("access_key.%d", access.UserID), time.Duration(access.TTL)).Err()
	if err != nil {
		return err
	}
	err = d.cache.Set(context.Background(), refresh.Token, fmt.Sprintf("refresh_key.%d", refresh.UserID), time.Duration(refresh.TTL)).Err()
	return err
}
