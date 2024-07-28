package storage

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type CacheDB struct {
	db *redis.Client
}

func NewCacheDB(db *redis.Client) *CacheDB {
	return &CacheDB{db: db}
}

func (d *CacheDB) SetToken(token string, ttl time.Duration) error {
	var status *redis.StatusCmd = d.db.Set(d.db.Context(), "token", token, ttl)

	return status.Err()
}
