package storage

import (
	"context"
	"errors"
	"github.com/MarshMeg/auth-micro_service.git/src/types"
	"github.com/MarshMeg/auth-micro_service.git/src/types/user"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strings"
)

type UserStorage struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewUserStorage(db *gorm.DB, redis *redis.Client) *UserStorage {
	return &UserStorage{
		db:    db,
		cache: redis,
	}
}

func (d *UserStorage) CreateUser(user user.User) (user.User, error) {
	result := d.db.Create(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (d *UserStorage) GetUserIDByToken(t string) (string, int, error) {
	res, err := d.cache.Get(context.Background(), t).Result()
	if err != nil {
		return "", 0, err
	}
	parts := strings.Split(res, "@")
	if len(parts) != 2 {
		return "", 0, errors.New("incorrect key")
	}
	return parts[0], types.StrToInt(parts[1]), err
}

func (d *UserStorage) UpdateUser(user user.User) error {
	return d.db.Model(&user).Update("username", user.Username).Update("password", user.Password).Error
}

func (d *UserStorage) GetUsers(filters *user.User) ([]user.User, error) {
	var users []user.User
	res := *d.db.Select([]string{"id", "username", "role"})
	if filters.Id != 0 {
		res.Where("`id`=?", filters.Id)
	}
	if filters.Username != "" {
		res.Where("INSTR(username, ?) > 0", filters.Username)
	}
	if filters.RoleName != "" {
		res.Where("`role_name`=?", filters.RoleName)
	}
	res.Find(&users)
	return users, res.Error
}

func (d *UserStorage) GetUser(filters *user.User) (user.User, error) {
	users, err := d.GetUsers(filters)
	return users[0], err
}

func (d *UserStorage) DeleteUser(user *user.User) error {
	return d.db.Delete(user).Error
}
