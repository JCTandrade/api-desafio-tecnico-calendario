package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/api"
)

const docsPage = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Holidays API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>window.ui = SwaggerUIBundle({url: "/openapi.yaml", dom_id: "#swagger-ui"});</script>
</body>
</html>`

// NewRouter registers the application routes and middlewares:
//
//	GET /health         liveness probe
//	GET /api/holidays   filtered holidays (JSON or XML)
//	GET /openapi.yaml   OpenAPI 3 specification
//	GET /docs           Swagger UI for the specification
func NewRouter(handler *Handler, logger *slog.Logger) *gin.Engine {
	router := gin.New()
	router.Use(requestLogger(logger), gin.Recovery())
	router.GET("/health", handler.Health)
	router.GET("/api/holidays", handler.ListHolidays)
	router.GET("/openapi.yaml", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/yaml; charset=utf-8", api.OpenAPISpec)
	})
	router.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(docsPage))
	})
	return router
}

// requestLogger logs one structured line per request once it has been served.
func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("http request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
			slog.String("client_ip", c.ClientIP()),
		)
	}
}
