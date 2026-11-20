package domain

import "strings"

type HolidayFilter struct {
	Year        *int
	Month       *int
	Type        string
	Inalienable *bool
}

func (f HolidayFilter) Matches(h Holiday) bool {
	return f.matchesYear(h) && f.matchesMonth(h) && f.matchesType(h) && f.matchesInalienable(h)
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
