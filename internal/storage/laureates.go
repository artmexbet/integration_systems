package storage

import (
	"context"
	"ris/internal/domain"
)

func (s *Storage) AddLaureates(ctx context.Context, laureates []domain.Laureate) error {
	return s.postgres.AddLaureates(ctx, laureates)
}
