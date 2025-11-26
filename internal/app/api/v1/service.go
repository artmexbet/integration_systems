package v1

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"ris/pkg/postgres/queries"
)

// Service defines the interface for the Nobel Prize API service
type Service interface {
	// Stats
	GetStats(ctx context.Context) (*StatsResponse, error)
	GetLastUpdate(ctx context.Context) (*LastUpdateResponse, error)

	// Laureates
	ListLaureates(ctx context.Context, page, perPage int) (*LaureateListResponse, error)
	GetLaureate(ctx context.Context, id int32) (*LaureateResponse, error)
	CreateLaureate(ctx context.Context, req *CreateLaureateRequest) (*LaureateResponse, error)
	UpdateLaureate(ctx context.Context, id int32, req *UpdateLaureateRequest) (*LaureateResponse, error)
	DeleteLaureate(ctx context.Context, id int32) error

	// Prizes
	ListPrizes(ctx context.Context, page, perPage int) (*PrizeListResponse, error)
	GetPrize(ctx context.Context, id int32) (*PrizeResponse, error)
	GetPrizesByCategory(ctx context.Context, category string) ([]PrizeResponse, error)
	GetPrizesByYear(ctx context.Context, year int32) ([]PrizeResponse, error)
	CreatePrize(ctx context.Context, req *CreatePrizeRequest) (*PrizeResponse, error)
	UpdatePrize(ctx context.Context, id int32, req *UpdatePrizeRequest) (*PrizeResponse, error)
	DeletePrize(ctx context.Context, id int32) error
	GetCategories(ctx context.Context) (*CategoriesResponse, error)
}

// NobelService implements the Service interface
type NobelService struct {
	queries *queries.Queries
}

// NewNobelService creates a new NobelService instance
func NewNobelService(q *queries.Queries) *NobelService {
	return &NobelService{queries: q}
}

// GetStats returns statistics about the dataset
func (s *NobelService) GetStats(ctx context.Context) (*StatsResponse, error) {
	stats, err := s.queries.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}
	return &StatsResponse{
		LaureatesCount:  stats.LaureatesCount,
		PrizesCount:     stats.PrizesCount,
		CategoriesCount: stats.CategoriesCount,
	}, nil
}

// GetLastUpdate returns the last update timestamp
func (s *NobelService) GetLastUpdate(ctx context.Context) (*LastUpdateResponse, error) {
	lastUpdate, err := s.queries.GetLastUpdate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get last update: %w", err)
	}

	// Try to convert the result to time.Time
	var t time.Time
	switch v := lastUpdate.(type) {
	case time.Time:
		t = v
	case *time.Time:
		if v != nil {
			t = *v
		}
	default:
		t = time.Now()
	}

	return &LastUpdateResponse{LastUpdate: t}, nil
}

// ListLaureates returns a paginated list of laureates
func (s *NobelService) ListLaureates(ctx context.Context, page, perPage int) (*LaureateListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	laureates, err := s.queries.ListLaureatesPaginated(ctx, queries.ListLaureatesPaginatedParams{
		Limit:  int32(perPage),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list laureates: %w", err)
	}

	total, err := s.queries.CountLaureates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count laureates: %w", err)
	}

	data := make([]LaureateResponse, len(laureates))
	for i, l := range laureates {
		data[i] = laureateToResponse(l)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &LaureateListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// GetLaureate returns a single laureate by ID
func (s *NobelService) GetLaureate(ctx context.Context, id int32) (*LaureateResponse, error) {
	laureate, err := s.queries.GetLaureate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("laureate not found: %w", err)
	}
	resp := laureateToResponse(laureate)
	return &resp, nil
}

// CreateLaureate creates a new laureate
func (s *NobelService) CreateLaureate(ctx context.Context, req *CreateLaureateRequest) (*LaureateResponse, error) {
	surname := pgtype.Text{}
	if req.Surname != "" {
		// pgtype.Text.Scan() only fails with nil input or non-string types,
		// which is not possible here since req.Surname is a string
		_ = surname.Scan(req.Surname)
	}

	laureate, err := s.queries.CreateLaureateSingle(ctx, queries.CreateLaureateSingleParams{
		ID:         req.ID,
		Firstname:  req.Firstname,
		Surname:    surname,
		Motivation: req.Motivation,
		Share:      req.Share,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create laureate: %w", err)
	}
	resp := laureateToResponse(laureate)
	return &resp, nil
}

// UpdateLaureate updates an existing laureate
func (s *NobelService) UpdateLaureate(ctx context.Context, id int32, req *UpdateLaureateRequest) (*LaureateResponse, error) {
	surname := pgtype.Text{}
	if req.Surname != "" {
		// pgtype.Text.Scan() only fails with nil input or non-string types,
		// which is not possible here since req.Surname is a string
		_ = surname.Scan(req.Surname)
	}

	laureate, err := s.queries.UpdateLaureate(ctx, queries.UpdateLaureateParams{
		ID:         id,
		Firstname:  req.Firstname,
		Surname:    surname,
		Motivation: req.Motivation,
		Share:      req.Share,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update laureate: %w", err)
	}
	resp := laureateToResponse(laureate)
	return &resp, nil
}

// DeleteLaureate deletes a laureate by ID
func (s *NobelService) DeleteLaureate(ctx context.Context, id int32) error {
	err := s.queries.DeleteLaureate(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete laureate: %w", err)
	}
	return nil
}

// ListPrizes returns a paginated list of prizes
func (s *NobelService) ListPrizes(ctx context.Context, page, perPage int) (*PrizeListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	prizes, err := s.queries.PrizesListPaginated(ctx, queries.PrizesListPaginatedParams{
		Limit:  int32(perPage),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list prizes: %w", err)
	}

	total, err := s.queries.CountPrizes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count prizes: %w", err)
	}

	data := make([]PrizeResponse, len(prizes))
	for i, p := range prizes {
		data[i] = prizeToResponse(p)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &PrizeListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// GetPrize returns a single prize with its laureates
func (s *NobelService) GetPrize(ctx context.Context, id int32) (*PrizeResponse, error) {
	prize, err := s.queries.GetPrize(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("prize not found: %w", err)
	}

	laureates, err := s.queries.GetLaureatesByPrizeId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get laureates: %w", err)
	}

	resp := prizeToResponse(prize)
	resp.Laureates = make([]LaureateResponse, len(laureates))
	for i, l := range laureates {
		resp.Laureates[i] = laureateToResponse(l)
	}

	return &resp, nil
}

// GetPrizesByCategory returns prizes filtered by category
func (s *NobelService) GetPrizesByCategory(ctx context.Context, category string) ([]PrizeResponse, error) {
	rows, err := s.queries.GetPrizesByCategoryWithLaureates(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get prizes by category: %w", err)
	}

	return aggregatePrizesWithLaureates(rows), nil
}

// GetPrizesByYear returns prizes filtered by year
func (s *NobelService) GetPrizesByYear(ctx context.Context, year int32) ([]PrizeResponse, error) {
	prizes, err := s.queries.PrizesByYear(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get prizes by year: %w", err)
	}

	data := make([]PrizeResponse, len(prizes))
	for i, p := range prizes {
		data[i] = prizeToResponse(p)
	}
	return data, nil
}

// CreatePrize creates a new prize
func (s *NobelService) CreatePrize(ctx context.Context, req *CreatePrizeRequest) (*PrizeResponse, error) {
	prize, err := s.queries.AddPrizeSingle(ctx, queries.AddPrizeSingleParams{
		Year:     req.Year,
		Category: req.Category,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create prize: %w", err)
	}

	// Link laureates if provided
	for _, laureateID := range req.LaureateIDs {
		err = s.queries.LinkLaureateToPrizeSingle(ctx, queries.LinkLaureateToPrizeSingleParams{
			PrizeID:    prize.ID,
			LaureateID: laureateID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to link laureate %d to prize: %w", laureateID, err)
		}
	}

	return s.GetPrize(ctx, prize.ID)
}

// UpdatePrize updates an existing prize
func (s *NobelService) UpdatePrize(ctx context.Context, id int32, req *UpdatePrizeRequest) (*PrizeResponse, error) {
	prize, err := s.queries.UpdatePrize(ctx, queries.UpdatePrizeParams{
		ID:       id,
		Year:     req.Year,
		Category: req.Category,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update prize: %w", err)
	}
	resp := prizeToResponse(prize)
	return &resp, nil
}

// DeletePrize deletes a prize by ID
func (s *NobelService) DeletePrize(ctx context.Context, id int32) error {
	err := s.queries.DeletePrize(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete prize: %w", err)
	}
	return nil
}

// GetCategories returns all unique prize categories
func (s *NobelService) GetCategories(ctx context.Context) (*CategoriesResponse, error) {
	categories, err := s.queries.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	if categories == nil {
		categories = []string{}
	}
	return &CategoriesResponse{Categories: categories}, nil
}

// Helper functions

func laureateToResponse(l queries.Laureate) LaureateResponse {
	resp := LaureateResponse{
		ID:         l.ID,
		Firstname:  l.Firstname,
		Motivation: l.Motivation,
		Share:      l.Share,
	}
	if l.Surname.Valid {
		resp.Surname = l.Surname.String
	}
	if l.UpdatedAt.Valid {
		t := l.UpdatedAt.Time.Format(time.RFC3339)
		resp.UpdatedAt = &t
	}
	return resp
}

func prizeToResponse(p queries.Prize) PrizeResponse {
	resp := PrizeResponse{
		ID:       p.ID,
		Year:     p.Year,
		Category: p.Category,
	}
	if p.UpdatedAt.Valid {
		t := p.UpdatedAt.Time.Format(time.RFC3339)
		resp.UpdatedAt = &t
	}
	return resp
}

func aggregatePrizesWithLaureates(rows []queries.GetPrizesByCategoryWithLaureatesRow) []PrizeResponse {
	if len(rows) == 0 {
		return []PrizeResponse{}
	}

	prizesMap := make(map[int32]*PrizeResponse)
	var orderedIDs []int32

	for _, row := range rows {
		prize, exists := prizesMap[row.PrizeID]
		if !exists {
			prize = &PrizeResponse{
				ID:        row.PrizeID,
				Year:      row.Year,
				Category:  row.Category,
				Laureates: []LaureateResponse{},
			}
			prizesMap[row.PrizeID] = prize
			orderedIDs = append(orderedIDs, row.PrizeID)
		}

		if row.LaureateID.Valid {
			laureate := LaureateResponse{
				ID: row.LaureateID.Int32,
			}
			if row.Firstname.Valid {
				laureate.Firstname = row.Firstname.String
			}
			if row.Surname.Valid {
				laureate.Surname = row.Surname.String
			}
			if row.Motivation.Valid {
				laureate.Motivation = row.Motivation.String
			}
			if row.Share.Valid {
				laureate.Share = row.Share.Int32
			}
			prize.Laureates = append(prize.Laureates, laureate)
		}
	}

	result := make([]PrizeResponse, len(orderedIDs))
	for i, id := range orderedIDs {
		result[i] = *prizesMap[id]
	}
	return result
}
