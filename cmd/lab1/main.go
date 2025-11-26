package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"ris/internal/parser"
	"ris/internal/storage"
	"time"
)

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	handler := slog.NewTextHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/ris")
	if err != nil {
		panic(err)
	}
	store := storage.NewStorage(pool)

	p := parser.NewParser(parser.Config{Url: "http://api.nobelprize.org/v1/prize.json"}, store)
	err = p.ParseAndStore(ctx)
	if err != nil {
		slog.Error("Error parsing and store ", "err", err)
	}
}
