package repository

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
)

type countingSource struct {
	calls    int
	holidays []domain.Holiday
	err      error
}

func (s *countingSource) GetAll(context.Context) ([]domain.Holiday, error) {
	s.calls++
	return s.holidays, s.err
}

var discardLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestInMemoryRepositoryLoadsSourceOnlyOnce(t *testing.T) {
	source := &countingSource{holidays: []domain.Holiday{{Title: "Año Nuevo", Date: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}}}

	repo, err := NewInMemoryHolidayRepository(context.Background(), source, discardLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for range 3 {
		got, err := repo.GetAll(context.Background())
		if err != nil || len(got) != 1 {
			t.Fatalf("unexpected result: %+v, %v", got, err)
		}
	}

	if source.calls != 1 {
		t.Errorf("source called %d times, want 1", source.calls)
	}
}

func TestInMemoryRepositoryReturnsCopy(t *testing.T) {
	source := &countingSource{holidays: []domain.Holiday{{Title: "Original"}}}
	repo, _ := NewInMemoryHolidayRepository(context.Background(), source, discardLogger)

	first, _ := repo.GetAll(context.Background())
	first[0].Title = "Mutated"
	second, _ := repo.GetAll(context.Background())

	if second[0].Title != "Original" {
		t.Errorf("cache was mutated: %q", second[0].Title)
	}
}

func TestInMemoryRepositoryFailsWhenSourceFails(t *testing.T) {
	cause := errors.New("boom")

	_, err := NewInMemoryHolidayRepository(context.Background(), &countingSource{err: cause}, discardLogger)

	if !errors.Is(err, cause) {
		t.Errorf("expected wrapped error, got %v", err)
	}
}
