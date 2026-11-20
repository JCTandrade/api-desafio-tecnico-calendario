package ports

import (
	"context"
	"errors"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
)

var ErrUpstream = errors.New("upstream holidays source failed")

type HolidayRepository interface {
	GetAll(ctx context.Context) ([]domain.Holiday, error)
}

type HolidayService interface {
	GetHolidays(ctx context.Context, filter domain.HolidayFilter) ([]domain.Holiday, error)
}
