package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain/ports"
)

type Handler struct {
	service ports.HolidayService
}

func NewHandler(service ports.HolidayService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListHolidays(w http.ResponseWriter, r *http.Request) {
	filter, err := parseFilter(r.URL.Query())
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	holidays, err := h.service.GetHolidays(r.Context(), filter)
	if err != nil {
		writeError(w, statusFor(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toHolidayResponses(holidays))
}

func parseFilter(query url.Values) (domain.HolidayFilter, error) {
	year, err := parseOptionalInt(query, "year", 1, 9999)
	if err != nil {
		return domain.HolidayFilter{}, err
	}
	month, err := parseOptionalInt(query, "month", 1, 12)
	if err != nil {
		return domain.HolidayFilter{}, err
	}
	inalienable, err := parseOptionalBool(query, "inalienable")
	if err != nil {
		return domain.HolidayFilter{}, err
	}
	return domain.HolidayFilter{Year: year, Month: month, Type: query.Get("type"), Inalienable: inalienable}, nil
}

func parseOptionalInt(query url.Values, key string, min, max int) (*int, error) {
	raw := query.Get(key)
	if raw == "" {
		return nil, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < min || value > max {
		return nil, fmt.Errorf("invalid %s %q: must be an integer between %d and %d", key, raw, min, max)
	}
	return &value, nil
}

func parseOptionalBool(query url.Values, key string) (*bool, error) {
	raw := query.Get(key)
	if raw == "" {
		return nil, nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid %s %q: must be true or false", key, raw)
	}
	return &value, nil
}

func statusFor(err error) int {
	if errors.Is(err, ports.ErrUpstream) {
		return http.StatusBadGateway
	}
	return http.StatusInternalServerError
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
