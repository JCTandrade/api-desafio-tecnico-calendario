package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain/ports"
)

// BoostrHolidayRepository reads holidays from the boostr.cl holidays JSON service.
type BoostrHolidayRepository struct {
	client *http.Client
	url    string
}

// NewBoostrHolidayRepository builds a repository that queries url using client.
func NewBoostrHolidayRepository(client *http.Client, url string) *BoostrHolidayRepository {
	return &BoostrHolidayRepository{client: client, url: url}
}

// GetAll fetches and maps every holiday; failures are wrapped with ports.ErrUpstream.
func (r *BoostrHolidayRepository) GetAll(ctx context.Context) ([]domain.Holiday, error) {
	payload, err := r.fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ports.ErrUpstream, err)
	}
	holidays, err := payload.toDomain()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ports.ErrUpstream, err)
	}
	return holidays, nil
}

func (r *BoostrHolidayRepository) fetch(ctx context.Context) (boostrResponse, error) {
	var payload boostrResponse
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, r.url, nil)
	if err != nil {
		return payload, fmt.Errorf("building request: %w", err)
	}
	response, err := r.client.Do(request)
	if err != nil {
		return payload, fmt.Errorf("performing request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return payload, fmt.Errorf("unexpected HTTP status %d", response.StatusCode)
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return payload, fmt.Errorf("decoding response: %w", err)
	}
	return payload, nil
}
