package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func connectDB() (*pgx.Conn, error) {

	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	databaseURL := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(
		context.Background(),
		databaseURL,
	)

	if err != nil {
		return nil, err
	}

	fmt.Println("Database connected successfully")

	return conn, nil
}