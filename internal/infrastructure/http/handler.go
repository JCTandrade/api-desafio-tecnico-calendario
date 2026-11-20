package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain"
	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/domain/ports"
)

const (
	formatJSON = "json"
	formatXML  = "xml"
)

// Handler adapts HTTP requests to the holidays use cases.
type Handler struct {
	service ports.HolidayService
	logger  *slog.Logger
}

// NewHandler builds a Handler that delegates to service and logs through logger.
func NewHandler(service ports.HolidayService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// Health reports that the process is up.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ListHolidays returns the holidays matching the query filters, as JSON or XML.
func (h *Handler) ListHolidays(c *gin.Context) {
	format, err := negotiateFormat(c)
	if err != nil {
		h.respondError(c, formatJSON, http.StatusBadRequest, err)
		return
	}
	filter, err := parseFilter(c)
	if err != nil {
		h.respondError(c, format, http.StatusBadRequest, err)
		return
	}
	holidays, err := h.service.GetHolidays(c.Request.Context(), filter)
	if err != nil {
		h.respondError(c, format, statusFor(err), err)
		return
	}
	responses := toHolidayResponses(holidays)
	if format == formatXML {
		c.XML(http.StatusOK, holidaysXML{Holidays: responses})
		return
	}
	c.JSON(http.StatusOK, responses)
}

// negotiateFormat picks the response format: the "format" query parameter wins over the
// Accept header, and JSON is the default.
func negotiateFormat(c *gin.Context) (string, error) {
	switch raw := strings.ToLower(c.Query("format")); raw {
	case "":
		if c.NegotiateFormat(gin.MIMEJSON, gin.MIMEXML) == gin.MIMEXML {
			return formatXML, nil
		}
		return formatJSON, nil
	case formatJSON, formatXML:
		return raw, nil
	default:
		return "", fmt.Errorf("invalid format %q: must be json or xml", raw)
	}
}

func parseFilter(c *gin.Context) (domain.HolidayFilter, error) {
	year, err := parseOptionalInt(c, "year", 1, 9999)
	if err != nil {
		return domain.HolidayFilter{}, err
	}
	month, err := parseOptionalInt(c, "month", 1, 12)
	if err != nil {
		return domain.HolidayFilter{}, err
	}
	inalienable, err := parseOptionalBool(c, "inalienable")
	if err != nil {
		return domain.HolidayFilter{}, err
	}
	from, err := parseOptionalDate(c, "from")
	if err != nil {
		return domain.HolidayFilter{}, err
	}
	to, err := parseOptionalDate(c, "to")
	if err != nil {
		return domain.HolidayFilter{}, err
	}
	return domain.HolidayFilter{
		Year:        year,
		Month:       month,
		Type:        strings.TrimSpace(c.Query("type")),
		Inalienable: inalienable,
		From:        from,
		To:          to,
	}, nil
}

func parseOptionalInt(c *gin.Context, key string, min, max int) (*int, error) {
	raw := c.Query(key)
	if raw == "" {
		return nil, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < min || value > max {
		return nil, fmt.Errorf("invalid %s %q: must be an integer between %d and %d", key, raw, min, max)
	}
	return &value, nil
}

func parseOptionalBool(c *gin.Context, key string) (*bool, error) {
	raw := c.Query(key)
	if raw == "" {
		return nil, nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid %s %q: must be true or false", key, raw)
	}
	return &value, nil
}

func parseOptionalDate(c *gin.Context, key string) (*time.Time, error) {
	raw := c.Query(key)
	if raw == "" {
		return nil, nil
	}
	value, err := time.Parse(dateLayout, raw)
	if err != nil {
		return nil, fmt.Errorf("invalid %s %q: must use the YYYY-MM-DD format", key, raw)
	}
	return &value, nil
}

func statusFor(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidDateRange):
		return http.StatusBadRequest
	case errors.Is(err, ports.ErrUpstream):
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

func (h *Handler) respondError(c *gin.Context, format string, status int, err error) {
	if status >= http.StatusInternalServerError {
		h.logger.Error("request failed", slog.String("path", c.Request.URL.Path), slog.Any("error", err))
	} else {
		h.logger.Warn("invalid request", slog.String("path", c.Request.URL.Path), slog.Any("error", err))
	}
	payload := errorResponse{Error: err.Error()}
	if format == formatXML {
		c.XML(status, payload)
		return
	}
	c.JSON(status, payload)
}
