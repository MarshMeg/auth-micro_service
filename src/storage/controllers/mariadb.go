package controllers

import (
	"github.com/DikosAs/auth-micro_service.git/src/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlClient(link string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(link), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&types.User{})
	if err != nil {
		return nil, err
	}

	return db, err
}
