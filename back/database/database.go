package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func ConnectDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
		dbURL = os.Getenv("DATABASE_URL")
	}
	var err error
	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return err
	}
	return nil
}

func CloseDB() {
	DB.Close()
}
