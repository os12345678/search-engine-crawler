package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectPGX() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig("postgres://admin:admin@localhost:5432/postgres")
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return db, nil
}