package storage

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"ris/pkg/postgres"
)

type Storage struct {
	postgres *postgres.Postgres
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		postgres: postgres.NewPostgres(pool),
	}
}
