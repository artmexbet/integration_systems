package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"ris/pkg/postgres/queries"
)

type Postgres struct {
	pool *pgxpool.Pool
	q    *queries.Queries
}

func NewPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{
		pool: pool,
		q:    queries.New(pool),
	}
}
