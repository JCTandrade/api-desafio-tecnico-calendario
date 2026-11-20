// Package ports declares the contracts between application and infrastructure.
package ports

import (
	"context"
	"errors"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
)

// ErrUpstream signals that the external holidays source failed or returned invalid data.
var ErrUpstream = errors.New("upstream holidays source failed")

// HolidayRepository is the driven (outbound) port that provides all known holidays.
type HolidayRepository interface {
	GetAll(ctx context.Context) ([]domain.Holiday, error)
}

// HolidayService is the driving (inbound) port that exposes the holidays use cases.
type HolidayService interface {
	// GetHolidays returns the holidays matching the filter, ordered by date.
	GetHolidays(ctx context.Context, filter domain.HolidayFilter) ([]domain.Holiday, error)
}
