// Package repository contains the adapters that fetch holidays from external sources.
package repository

import (
	"fmt"
	"time"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
)

const (
	successStatus = "success"
	dateLayout    = "2006-01-02"
)

type boostrResponse struct {
	Status string          `json:"status"`
	Data   []boostrHoliday `json:"data"`
}

type boostrHoliday struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Inalienable bool   `json:"inalienable"`
	Extra       string `json:"extra"`
}

func (r boostrResponse) toDomain() ([]domain.Holiday, error) {
	if r.Status != successStatus {
		return nil, fmt.Errorf("unexpected status %q", r.Status)
	}
	holidays := make([]domain.Holiday, 0, len(r.Data))
	for _, item := range r.Data {
		holiday, err := item.toDomain()
		if err != nil {
			return nil, err
		}
		holidays = append(holidays, holiday)
	}
	return holidays, nil
}

func (h boostrHoliday) toDomain() (domain.Holiday, error) {
	date, err := time.Parse(dateLayout, h.Date)
	if err != nil {
		return domain.Holiday{}, fmt.Errorf("parsing date %q: %w", h.Date, err)
	}
	return domain.Holiday{
		Date:        date,
		Title:       h.Title,
		Type:        h.Type,
		Inalienable: h.Inalienable,
		Description: h.Extra,
	}, nil
}
