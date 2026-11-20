package domain

import (
	"errors"
	"strings"
	"time"
)

// ErrInvalidDateRange is returned when the range start is after its end.
var ErrInvalidDateRange = errors.New("invalid date range: 'from' must not be after 'to'")

// HolidayFilter narrows a list of holidays. A zero-valued field means "no restriction".
type HolidayFilter struct {
	// Year restricts to a single calendar year.
	Year *int
	// Month restricts to a single month (1-12).
	Month *int
	// Type restricts to a holiday type (case-insensitive), e.g. "Religioso".
	Type string
	// Inalienable restricts by the inalienable flag.
	Inalienable *bool
	// From is the inclusive start of a date range.
	From *time.Time
	// To is the inclusive end of a date range.
	To *time.Time
}

// Validate checks the filter is internally consistent.
func (f HolidayFilter) Validate() error {
	if f.From != nil && f.To != nil && f.From.After(*f.To) {
		return ErrInvalidDateRange
	}
	return nil
}

// Matches reports whether the holiday satisfies every restriction of the filter.
func (f HolidayFilter) Matches(h Holiday) bool {
	return f.matchesYear(h) && f.matchesMonth(h) && f.matchesType(h) &&
		f.matchesInalienable(h) && f.matchesRange(h)
}

func (f HolidayFilter) matchesYear(h Holiday) bool {
	return f.Year == nil || h.Date.Year() == *f.Year
}

func (f HolidayFilter) matchesMonth(h Holiday) bool {
	return f.Month == nil || int(h.Date.Month()) == *f.Month
}

func (f HolidayFilter) matchesType(h Holiday) bool {
	return f.Type == "" || strings.EqualFold(f.Type, h.Type)
}

func (f HolidayFilter) matchesInalienable(h Holiday) bool {
	return f.Inalienable == nil || h.Inalienable == *f.Inalienable
}

func (f HolidayFilter) matchesRange(h Holiday) bool {
	if f.From != nil && h.Date.Before(*f.From) {
		return false
	}
	return f.To == nil || !h.Date.After(*f.To)
}
