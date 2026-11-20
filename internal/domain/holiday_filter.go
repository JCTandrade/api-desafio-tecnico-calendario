package domain

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidDateRange = errors.New("invalid date range: 'from' must not be after 'to'")

type HolidayFilter struct {
	Year *int

	Month *int

	Type string

	Inalienable *bool

	From *time.Time

	To *time.Time
}

func (f HolidayFilter) Validate() error {
	if f.From != nil && f.To != nil && f.From.After(*f.To) {
		return ErrInvalidDateRange
	}
	return nil
}

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
