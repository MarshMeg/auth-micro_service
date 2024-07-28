package controllers

import (
	"github.com/go-redis/redis/v8"
)

type AccessToken struct {
	UserID int    `json:"user_id"`
	Token  string `json:"access_token"`
	TTL    int    `json:"ttl"`
}

type RefreshToken struct {
	UserID int    `json:"user_id"`
	Token  string `json:"refresh_token"`
	TTL    int    `json:"ttl"`
}

func NewRedisClient(cfg *redis.Options) *redis.Client {
	db := redis.NewClient(cfg)
	return db
}
