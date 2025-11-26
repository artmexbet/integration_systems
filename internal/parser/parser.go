package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"ris/internal/domain"
)

type storage interface {
	AddLaureates(context.Context, []domain.Laureate) error
	AddPrizes([]domain.Prize) ([]int32, error)
	LinkLaureatesToPrizes(ctx context.Context, prizeId int32, laureates []domain.Laureate) error
}

type Config struct {
	Url string
}

type Parser struct {
	cfg     Config
	storage storage
	client  http.Client
}

func NewParser(cfg Config, storage storage) *Parser {
	return &Parser{
		cfg:     cfg,
		storage: storage,
		client:  http.Client{},
	}
}

func (p *Parser) ParseAndStore(ctx context.Context) error {
	// Fetch data from the URL
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.cfg.Url, nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	slog.Info("Request URL: ", req.URL.String())

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()
	slog.Info("Response Status: ", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	var nobelResponse domain.NobelResponse
	if err := json.NewDecoder(resp.Body).Decode(&nobelResponse); err != nil {
		return fmt.Errorf("could not decode response body: %w", err)
	}
	// Transform raw data to domain models
	prizes := make([]domain.Prize, 0, len(nobelResponse.Prizes))
	laureateMap := make(map[int32]domain.Laureate)
	for _, rawPrize := range nobelResponse.Prizes {
		prize := rawPrize.ToPrize()
		prizes = append(prizes, prize)
		for _, laureate := range prize.Laureates {
			laureateMap[laureate.Id] = laureate
		}
	}
	slog.Info("Inserted prizes", "prizes", prizes)

	laureates := make([]domain.Laureate, 0, len(laureateMap))
	for _, laureate := range laureateMap {
		laureates = append(laureates, laureate)
	}
	err = p.storage.AddLaureates(ctx, laureates)
	if err != nil {
		return fmt.Errorf("could not add laureates: %w", err)
	}

	prizesIds, err := p.storage.AddPrizes(prizes)
	if err != nil {
		return fmt.Errorf("could not add prizes: %w", err)
	}

	for i, prize := range prizes {
		err = p.storage.LinkLaureatesToPrizes(ctx, prizesIds[i], prize.Laureates)
		if err != nil {
			return fmt.Errorf("could not link laureates to prize: %w", err)
		}
	}
	return nil
}
