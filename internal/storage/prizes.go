package storage

import (
	"context"
	"ris/internal/domain"
)

func (s *Storage) AddPrizes(prizes []domain.Prize) ([]int32, error) {
	return s.postgres.AddPrizes(prizes)
}

func (s *Storage) LinkLaureatesToPrizes(ctx context.Context, prizeId int32, laureates []domain.Laureate) error {
	return s.postgres.LinkLaureatesToPrizes(ctx, prizeId, laureates)
}
