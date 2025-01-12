package db

import (
	"database/sql"
	"log"
)

func NewSQLStorage(cfg string) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitStorage(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}
	log.Println("db-svc connected")
	return nil
}
