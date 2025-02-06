package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	connStr := "postgres://postgres:password@localhost:5432/postgres-demo?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		slog.Error("Failed to initialise connection to database", "error", err)
		return
	}
	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		slog.Error("Failed to ping the database", "error", err)
		return
	}

	slog.Info("Database connection initialised successfully!")
}
