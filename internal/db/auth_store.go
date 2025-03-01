package db

import "database/sql"

type User struct {
	Id       string `json:"id"`
	Provider string `json:"provider"`
}

type AuthStore interface {
}

type SqlAuthStore struct {
	db *sql.DB
}

func NewAuthStore(db *sql.DB) *SqlAuthStore {
	return &SqlAuthStore{db: db}
}
