package repository

import (
	"github.com/jmoiron/sqlx"
	"strings"
)

func NewMysqlDB(link string) (*sqlx.DB, error) {
	db, err := sqlx.Connect(strings.Split(link, "://")[0], strings.Split(link, "://")[1])
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
