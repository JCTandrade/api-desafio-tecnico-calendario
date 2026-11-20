package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain/ports"
)

// InMemoryHolidayRepository serves holidays from memory. The data is loaded exactly once,
// when the repository is created, so the upstream source is never called per request.
type InMemoryHolidayRepository struct {
	holidays []domain.Holiday
}

// NewInMemoryHolidayRepository eagerly loads every holiday from source.
// It fails if the source cannot be read, so the application never starts without data.
func NewInMemoryHolidayRepository(ctx context.Context, source ports.HolidayRepository, logger *slog.Logger) (*InMemoryHolidayRepository, error) {
	holidays, err := source.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading holidays at startup: %w", err)
	}
	logger.Info("holidays loaded into memory", slog.Int("count", len(holidays)))
	return &InMemoryHolidayRepository{holidays: holidays}, nil
}

// GetAll returns a copy of the cached holidays, so callers cannot mutate the cache.
func (r *InMemoryHolidayRepository) GetAll(_ context.Context) ([]domain.Holiday, error) {
	return append([]domain.Holiday(nil), r.holidays...), nil
}
