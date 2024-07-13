package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go_back/src/repository/objects"
	"time"
)

type Authorisation interface {
	CreateUser(user objects.User) (int, error)
	GetUser(username, password string) (objects.User, error)
}

type Cache interface {
	SetToken(token string, ttl time.Duration) error
}

type Repository struct {
	Authorisation
	Cache
}

func NewRepository(fileDB *sqlx.DB, cacheDB *redis.Client) *Repository {
	return &Repository{
		Authorisation: NewAuthDB(fileDB),
		Cache:         NewCacheDB(cacheDB),
	}
}
