// Package main wires the application together.
package main

import (
	"log"
	"net/http"

	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/application"
	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/infrastructure/config"
	httpadapter "github.com/JCTandrade/api-desafio-tecnico-calendario/internal/infrastructure/http"
	"github.com/JCTandrade/api-desafio-tecnico-calendario/internal/infrastructure/repository"
)

func main() {
	settings := config.Load()
	client := &http.Client{Timeout: settings.HTTPTimeout}
	repo := repository.NewBoostrHolidayRepository(client, settings.HolidaysURL)
	service := application.NewHolidayService(repo)
	router := httpadapter.NewRouter(httpadapter.NewHandler(service), "web")

	log.Printf("listening on :%s", settings.Port)
	log.Fatal(http.ListenAndServe(":"+settings.Port, router))
}
