package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/application"
	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/infrastructure/config"
	httpadapter "github.com/JCTandrade/api-desafio-tecnico-calendario/internal/infrastructure/http"
	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/infrastructure/repository"
)

const shutdownTimeout = 10 * time.Second

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("application stopped", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	settings := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	source := repository.NewBoostrHolidayRepository(&http.Client{Timeout: settings.HTTPTimeout}, settings.HolidaysURL)
	repo, err := repository.NewInMemoryHolidayRepository(ctx, source, logger)
	if err != nil {
		return err
	}

	service := application.NewHolidayService(repo)
	router := httpadapter.NewRouter(httpadapter.NewHandler(service, logger), logger)
	server := &http.Server{Addr: ":" + settings.Port, Handler: router, ReadHeaderTimeout: 5 * time.Second}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server listening", slog.String("addr", server.Addr))
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		logger.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
