package http

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain/ports"
)

type fakeService struct {
	holidays []domain.Holiday
	err      error
}

func (f fakeService) GetHolidays(context.Context, domain.HolidayFilter) ([]domain.Holiday, error) {
	return f.holidays, f.err
}

func performRequest(service fakeService, target string, headers ...string) *httptest.ResponseRecorder {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(NewHandler(service, logger), logger)
	request := httptest.NewRequest(http.MethodGet, target, nil)
	for i := 0; i+1 < len(headers); i += 2 {
		request.Header.Set(headers[i], headers[i+1])
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestListHolidays(t *testing.T) {
	holiday := domain.Holiday{Date: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Title: "Año Nuevo"}
	tests := []struct {
		name    string
		service fakeService
		target  string
		want    int
	}{
		{"ok", fakeService{holidays: []domain.Holiday{holiday}}, "/api/holidays?year=2026", http.StatusOK},
		{"invalid month", fakeService{}, "/api/holidays?month=13", http.StatusBadRequest},
		{"invalid inalienable", fakeService{}, "/api/holidays?inalienable=x", http.StatusBadRequest},
		{"invalid from", fakeService{}, "/api/holidays?from=01-01-2026", http.StatusBadRequest},
		{"invalid to", fakeService{}, "/api/holidays?to=nope", http.StatusBadRequest},
		{"inverted range", fakeService{err: domain.ErrInvalidDateRange}, "/api/holidays?from=2026-06-01&to=2026-01-01", http.StatusBadRequest},
		{"invalid format", fakeService{}, "/api/holidays?format=csv", http.StatusBadRequest},
		{"upstream failure", fakeService{err: fmt.Errorf("wrap: %w", ports.ErrUpstream)}, "/api/holidays", http.StatusBadGateway},
		{"internal failure", fakeService{err: fmt.Errorf("other")}, "/api/holidays", http.StatusInternalServerError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := performRequest(test.service, test.target).Code; got != test.want {
				t.Errorf("status = %d, want %d", got, test.want)
			}
		})
	}
}

func TestListHolidaysContentNegotiation(t *testing.T) {
	service := fakeService{holidays: []domain.Holiday{{
		Date: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), Title: "Día del Trabajador",
		Type: "Civil", Inalienable: true, Extra: "Civil e Irrenunciable",
	}}}
	tests := []struct {
		name         string
		target       string
		headers      []string
		wantType     string
		wantContains []string
	}{
		{"json by default", "/api/holidays", nil, "application/json",
			[]string{`"date":"2026-05-01"`, `"title":"Día del Trabajador"`, `"type":"Civil"`, `"inalienable":true`, `"extra":"Civil e Irrenunciable"`}},
		{"xml by accept header", "/api/holidays", []string{"Accept", "application/xml"}, "application/xml",
			[]string{"<holidays>", "<holiday>", "<date>2026-05-01</date>", "<inalienable>true</inalienable>", "<extra>Civil e Irrenunciable</extra>"}},
		{"xml by format parameter", "/api/holidays?format=xml", nil, "application/xml",
			[]string{"<holidays>", "<title>Día del Trabajador</title>"}},
		{"format parameter wins over accept", "/api/holidays?format=json", []string{"Accept", "application/xml"}, "application/json",
			[]string{`"date":"2026-05-01"`}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := performRequest(service, test.target, test.headers...)

			if got := recorder.Header().Get("Content-Type"); !strings.HasPrefix(got, test.wantType) {
				t.Errorf("Content-Type = %q, want prefix %q", got, test.wantType)
			}
			for _, want := range test.wantContains {
				if !strings.Contains(recorder.Body.String(), want) {
					t.Errorf("body %q does not contain %q", recorder.Body.String(), want)
				}
			}
		})
	}
}

func TestErrorsAreNegotiatedToo(t *testing.T) {
	recorder := performRequest(fakeService{}, "/api/holidays?month=13", "Accept", "application/xml")

	if recorder.Code != http.StatusBadRequest || !strings.Contains(recorder.Body.String(), "<error>") {
		t.Errorf("unexpected response: %d %q", recorder.Code, recorder.Body.String())
	}
}

func TestOpenAPISpecIsServed(t *testing.T) {
	recorder := performRequest(fakeService{}, "/openapi.yaml")

	if recorder.Code != http.StatusOK || !strings.HasPrefix(recorder.Body.String(), "openapi: 3.") {
		t.Errorf("unexpected response: %d %.40q", recorder.Code, recorder.Body.String())
	}
}

func TestHealth(t *testing.T) {
	recorder := performRequest(fakeService{}, "/health")

	if recorder.Code != http.StatusOK || strings.TrimSpace(recorder.Body.String()) != `{"status":"ok"}` {
		t.Errorf("unexpected response: %d %q", recorder.Code, recorder.Body.String())
	}
}
