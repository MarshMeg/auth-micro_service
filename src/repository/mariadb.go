package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewMysqlClient(link string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(link), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}

	return db, err
}
