package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/pkg/errors"
	"time"
)

type DB struct {
	DB     *sql.DB
	logger *logging.Logger
}

func NewDB(dsn string, logger *logging.Logger) (*DB, error) {
	var db DB
	d, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "can't create database connection")
	}
	err = d.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "can't open connection to database")
	}

	db.DB = d

	db.DB.SetMaxOpenConns(10)
	db.DB.SetMaxIdleConns(5)
	db.DB.SetConnMaxLifetime(time.Minute * 5) // Prevent long-lived connections

	db.logger = logger

	return &db, err
}

func (db *DB) Close() error {
	err := db.DB.Close()
	if err != nil {
		return errors.Wrap(err, "error closing database connection")
	}
	return nil
}
