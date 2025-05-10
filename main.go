package main

import (
	"context"
	"log/slog"
	"main/models"
	"os"

	"github.com/jackc/pgx/v5"
)

// the go part of this project will host the backend and has roles ranging from auth
// to saving notes, decks, and progress

// env:
// POSTGRES_URL: the url to the postgres database
// API_ADDR: the addr to run the API on

func main() {
	ctx := context.Background()
	logger := slog.Default()

	postgres_url := os.Getenv("POSTGRES_URL")
	conn, err := pgx.Connect(ctx, postgres_url)
	if err != nil {
		logger.Debug("postgres", "url", postgres_url)
		logger.Error("pgx connection failed (fatal)", "err", err.Error())
		return
	} else {
		logger.Debug("pgx connection success")
	}
	defer conn.Close(ctx)

	if err := models.Init(ctx, conn); err != nil {
		logger.Error("models init failed", "err", err.Error())
		return
	} else {
		logger.Debug("models init success")
	}

	if err := Serve(conn); err != nil {
		logger.Error("server", "err", err.Error())
		return
	}
}
