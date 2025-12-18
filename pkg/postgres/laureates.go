package postgres

import (
	"context"
	"fmt"
	"ris/internal/domain"
	"ris/pkg/postgres/queries"

	"github.com/jackc/pgx/v5/pgtype"
)

func (p *Postgres) AddLaureates(ctx context.Context, laureates []domain.Laureate) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	params := make([]queries.CreateLaureateParams, 0, len(laureates))
	for _, laureate := range laureates {
		surname := pgtype.Text{}
		_ = surname.Scan(laureate.Surname)
		params = append(params, queries.CreateLaureateParams{
			ID:         laureate.Id,
			Firstname:  laureate.Firstname,
			Surname:    surname,
			Motivation: laureate.Motivation,
			Share:      laureate.Share,
		})
	}
	res := p.q.WithTx(tx).CreateLaureate(ctx, params)
	ok := true
	res.QueryRow(func(_ int, _ queries.Laureate, err error) {
		if err != nil {
			ok = false
		}
	})
	if !ok {
		return fmt.Errorf("could not create laureates")
	}
	_ = tx.Commit(ctx)
	return nil
}

func (p *Postgres) LinkLaureatesToPrizes(ctx context.Context, prizeId int32, laureates []domain.Laureate) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	params := make([]queries.LinkLaureateToPrizeParams, 0, len(laureates))
	for _, laureate := range laureates {
		params = append(params, queries.LinkLaureateToPrizeParams{
			LaureateID: laureate.Id,
			PrizeID:    prizeId,
		})
	}
	res := p.q.WithTx(tx).LinkLaureateToPrize(ctx, params)
	ok := true
	res.Exec(func(i int, err error) {
		if err != nil {
			ok = false
			return
		}
	})
	if !ok {
		return fmt.Errorf("could not link laureates to prize")
	}
	_ = tx.Commit(ctx)
	return nil
}
