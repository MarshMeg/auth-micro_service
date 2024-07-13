package repository

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DBNum    int
}

func NewRedisDB(cnf *RedisConfig) *redis.Client {
	var db *redis.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
		Password: cnf.Password,
		DB:       cnf.DBNum,
	})

	return db
}
