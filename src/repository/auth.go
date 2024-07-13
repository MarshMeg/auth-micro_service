package repository

import (
	"github.com/jmoiron/sqlx"
	"go_back/src/repository/objects"
)

type AuthDB struct {
	db *sqlx.DB
}

func NewAuthDB(db *sqlx.DB) *AuthDB {
	return &AuthDB{db: db}
}

func (d *AuthDB) CreateUser(user objects.User) (int, error) {
	result, err := d.db.Exec("INSERT INTO `Users`(`username`, `password`) VALUES(?, ?)", user.Username, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (d *AuthDB) GetUser(username, password string) (objects.User, error) {
	var user objects.User

	err := d.db.Get(&user, "SELECT `id` FROM `Users` WHERE `username`=? AND `password`=?", username, password)
	return user, err
}
