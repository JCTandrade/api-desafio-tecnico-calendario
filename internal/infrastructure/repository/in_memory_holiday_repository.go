package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain/ports"
)

type InMemoryHolidayRepository struct {
	holidays []domain.Holiday
}

func NewInMemoryHolidayRepository(ctx context.Context, source ports.HolidayRepository, logger *slog.Logger) (*InMemoryHolidayRepository, error) {
	holidays, err := source.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading holidays at startup: %w", err)
	}
	logger.Info("holidays loaded into memory", slog.Int("count", len(holidays)))
	return &InMemoryHolidayRepository{holidays: holidays}, nil
}

func (r *InMemoryHolidayRepository) GetAll(_ context.Context) ([]domain.Holiday, error) {
	return append([]domain.Holiday(nil), r.holidays...), nil
}
