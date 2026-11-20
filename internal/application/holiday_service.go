package application

import (
	"context"
	"fmt"
	"sort"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain/ports"
)

type HolidayService struct {
	repository ports.HolidayRepository
}

func NewHolidayService(repository ports.HolidayRepository) *HolidayService {
	return &HolidayService{repository: repository}
}

func (s *HolidayService) GetHolidays(ctx context.Context, filter domain.HolidayFilter) ([]domain.Holiday, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}
	holidays, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting holidays from repository: %w", err)
	}
	filtered := make([]domain.Holiday, 0, len(holidays))
	for _, holiday := range holidays {
		if filter.Matches(holiday) {
			filtered = append(filtered, holiday)
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool { return filtered[i].Date.Before(filtered[j].Date) })
	return filtered, nil
}
