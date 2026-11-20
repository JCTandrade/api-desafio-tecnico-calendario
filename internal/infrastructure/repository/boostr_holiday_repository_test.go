package repository

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain/ports"
)

func newServer(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(body))
	}))
}

func TestGetAllSuccess(t *testing.T) {
	server := newServer(http.StatusOK, `{"status":"success","data":[{"date":"2026-01-01","title":"Año Nuevo","type":"Civil","inalienable":true,"extra":"Civil e Irrenunciable"}]}`)
	defer server.Close()

	holidays, err := NewBoostrHolidayRepository(server.Client(), server.URL).GetAll(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(holidays) != 1 || holidays[0].Title != "Año Nuevo" || holidays[0].Extra != "Civil e Irrenunciable" || !holidays[0].Inalienable {
		t.Errorf("unexpected holidays: %+v", holidays)
	}
}

func TestGetAllFailures(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
	}{
		{"status not success", http.StatusOK, `{"status":"error","data":[]}`},
		{"invalid json", http.StatusOK, `not json`},
		{"http 500", http.StatusInternalServerError, ``},
		{"invalid date", http.StatusOK, `{"status":"success","data":[{"date":"bad"}]}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := newServer(test.statusCode, test.body)
			defer server.Close()

			_, err := NewBoostrHolidayRepository(server.Client(), server.URL).GetAll(context.Background())

			if !errors.Is(err, ports.ErrUpstream) {
				t.Errorf("expected upstream error, got %v", err)
			}
		})
	}
}
