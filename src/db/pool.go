package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

// get a connection pool to the database with connection info via .env file

func ConnectPGX() (*pgxpool.Pool, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	user:= os.Getenv("DB_USER")
	host:= os.Getenv("DB_HOST")
	dbname:= os.Getenv("DB_NAME")
	password:= os.Getenv("DB_PASSWORD")
	port:= os.Getenv("DB_PORT")


	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	// Connect to the database
	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return db, nil
}