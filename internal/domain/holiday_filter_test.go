package domain

import (
	"errors"
	"testing"
	"time"
)

func TestHolidayFilterMatches(t *testing.T) {
	holiday := Holiday{
		Date:        time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		Type:        "Civil",
		Inalienable: true,
	}
	year, otherYear := 2026, 2025
	month, otherMonth := 1, 2
	yes, no := true, false
	day := holiday.Date
	before, after := day.AddDate(0, 0, -1), day.AddDate(0, 0, 1)

	tests := []struct {
		name   string
		filter HolidayFilter
		want   bool
	}{
		{"empty filter", HolidayFilter{}, true},
		{"year match", HolidayFilter{Year: &year}, true},
		{"year mismatch", HolidayFilter{Year: &otherYear}, false},
		{"month match", HolidayFilter{Month: &month}, true},
		{"month mismatch", HolidayFilter{Month: &otherMonth}, false},
		{"type case insensitive", HolidayFilter{Type: "civil"}, true},
		{"type mismatch", HolidayFilter{Type: "Religioso"}, false},
		{"inalienable match", HolidayFilter{Inalienable: &yes}, true},
		{"inalienable mismatch", HolidayFilter{Inalienable: &no}, false},
		{"combined", HolidayFilter{Year: &year, Month: &month, Type: "CIVIL", Inalienable: &yes}, true},
		{"range contains date", HolidayFilter{From: &before, To: &after}, true},
		{"range bounds are inclusive", HolidayFilter{From: &day, To: &day}, true},
		{"before range start", HolidayFilter{From: &after}, false},
		{"after range end", HolidayFilter{To: &before}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.filter.Matches(holiday); got != test.want {
				t.Errorf("Matches() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestHolidayFilterValidate(t *testing.T) {
	early := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	late := early.AddDate(0, 1, 0)

	if err := (HolidayFilter{From: &early, To: &late}).Validate(); err != nil {
		t.Errorf("valid range rejected: %v", err)
	}
	if err := (HolidayFilter{From: &late, To: &early}).Validate(); !errors.Is(err, ErrInvalidDateRange) {
		t.Errorf("expected ErrInvalidDateRange, got %v", err)
	}
}
