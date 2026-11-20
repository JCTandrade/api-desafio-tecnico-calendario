package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
)

type fakeRepository struct {
	holidays []domain.Holiday
	err      error
}

func (f fakeRepository) GetAll(context.Context) ([]domain.Holiday, error) {
	return f.holidays, f.err
}

func TestGetHolidaysFiltersResults(t *testing.T) {
	repository := fakeRepository{holidays: []domain.Holiday{
		{Title: "Año Nuevo", Type: "Civil", Date: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{Title: "Viernes Santo", Type: "Religioso", Date: time.Date(2026, time.April, 3, 0, 0, 0, 0, time.UTC)},
	}}
	month := 4

	got, err := NewHolidayService(repository).GetHolidays(context.Background(), domain.HolidayFilter{Month: &month})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Title != "Viernes Santo" {
		t.Errorf("unexpected holidays: %+v", got)
	}
}

func TestGetHolidaysWrapsRepositoryError(t *testing.T) {
	cause := errors.New("boom")

	_, err := NewHolidayService(fakeRepository{err: cause}).GetHolidays(context.Background(), domain.HolidayFilter{})

	if !errors.Is(err, cause) {
		t.Errorf("expected wrapped error, got %v", err)
	}
}

func TestGetHolidaysOrdersByDateAndFiltersByRange(t *testing.T) {
	date := func(month time.Month, day int) time.Time { return time.Date(2026, month, day, 0, 0, 0, 0, time.UTC) }
	repository := fakeRepository{holidays: []domain.Holiday{
		{Title: "C", Date: date(time.December, 25)},
		{Title: "A", Date: date(time.March, 1)},
		{Title: "B", Date: date(time.May, 1)},
	}}
	from, to := date(time.February, 1), date(time.June, 30)

	got, err := NewHolidayService(repository).GetHolidays(context.Background(), domain.HolidayFilter{From: &from, To: &to})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].Title != "A" || got[1].Title != "B" {
		t.Errorf("unexpected holidays: %+v", got)
	}
}

func TestGetHolidaysRejectsInvertedRange(t *testing.T) {
	from := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, -1, 0)

	_, err := NewHolidayService(fakeRepository{}).GetHolidays(context.Background(), domain.HolidayFilter{From: &from, To: &to})

	if !errors.Is(err, domain.ErrInvalidDateRange) {
		t.Errorf("expected ErrInvalidDateRange, got %v", err)
	}
}
