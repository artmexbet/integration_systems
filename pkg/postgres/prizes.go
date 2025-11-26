package postgres

import (
	"context"
	"fmt"
	"ris/internal/domain"
	"ris/pkg/postgres/queries"
	"ris/pkg/utills"
)

func (p *Postgres) AddPrizes(prizes []domain.Prize) ([]int32, error) {
	ctx := context.Background()
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	params := make([]queries.AddPrizeParams, 0, len(prizes))
	for _, prize := range prizes {
		params = append(params, queries.AddPrizeParams{
			Year:     int32(utills.ParseStringToInt(prize.Year)),
			Category: prize.Category,
		})
	}

	prizesIds := make([]int32, len(prizes))
	res := p.q.WithTx(tx).AddPrize(ctx, params)
	ok := true
	res.QueryRow(func(i int, p queries.Prize, err error) {
		if err != nil {
			ok = false
			return
		}

		prizesIds[i] = p.ID
	})
	if !ok {
		return nil, fmt.Errorf("could not add prizes")
	}
	_ = tx.Commit(ctx)
	return prizesIds, nil
}
