package http

import "net/http"

func NewRouter(handler *Handler, webDir string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /api/holidays", handler.ListHolidays)
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	return mux
}
