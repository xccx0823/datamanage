package database

import (
	"github.com/jmoiron/sqlx"
)

type DbOption func(db *sqlx.DB)

// WithMaxOpenCons 设置最大连接数
func WithMaxOpenCons(n int) DbOption {
	return func(db *sqlx.DB) {
		db.SetMaxOpenConns(n)
	}
}

// WithMaxIdleCons 设置空闲连接数
func WithMaxIdleCons(n int) DbOption {
	return func(db *sqlx.DB) {
		db.SetMaxIdleConns(n)
	}
}
