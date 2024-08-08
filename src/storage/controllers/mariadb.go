package controllers

import (
	"github.com/MarshMeg/auth-micro_service.git/src/types/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlClient(link string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(link), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&user.User{})
	if err != nil {
		return nil, err
	}

	return db, err
}
