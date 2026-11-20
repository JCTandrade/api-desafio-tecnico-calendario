package http

import (
	"encoding/xml"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
)

const dateLayout = "2006-01-02"

type holidayResponse struct {
	XMLName     xml.Name `json:"-" xml:"holiday"`
	Date        string   `json:"date" xml:"date"`
	Title       string   `json:"title" xml:"title"`
	Type        string   `json:"type" xml:"type"`
	Inalienable bool     `json:"inalienable" xml:"inalienable"`
	Extra       string   `json:"extra" xml:"extra"`
}

type holidaysXML struct {
	XMLName  xml.Name          `xml:"holidays"`
	Holidays []holidayResponse `xml:"holiday"`
}

type errorResponse struct {
	XMLName xml.Name `json:"-" xml:"error"`
	Error   string   `json:"error" xml:"message"`
}

func toHolidayResponses(holidays []domain.Holiday) []holidayResponse {
	responses := make([]holidayResponse, 0, len(holidays))
	for _, holiday := range holidays {
		responses = append(responses, holidayResponse{
			Date:        holiday.Date.Format(dateLayout),
			Title:       holiday.Title,
			Type:        holiday.Type,
			Inalienable: holiday.Inalienable,
			Extra:       holiday.Extra,
		})
	}
	return responses
}
