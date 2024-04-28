package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
    err := godotenv.Load(".env")
    if err != nil {
        fmt.Println("Error loading .env file:", err)
    }
}

func createSQLPool() string {
    loadEnv() 

    dbUser := os.Getenv("DB_USER")
    dbHost := os.Getenv("DB_HOST")
    dbName := os.Getenv("DB_NAME")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbPort := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName,
    )

	return connectionString
}