// Package http exposes the holidays use cases over HTTP.
package http

import "github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"

const dateLayout = "2006-01-02"

type holidayResponse struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Inalienable bool   `json:"inalienable"`
	Description string `json:"description"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func toHolidayResponses(holidays []domain.Holiday) []holidayResponse {
	responses := make([]holidayResponse, 0, len(holidays))
	for _, holiday := range holidays {
		responses = append(responses, holidayResponse{
			Date:        holiday.Date.Format(dateLayout),
			Title:       holiday.Title,
			Type:        holiday.Type,
			Inalienable: holiday.Inalienable,
			Description: holiday.Description,
		})
	}
	return responses
}
