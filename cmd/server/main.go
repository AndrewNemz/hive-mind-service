package main

import (
	"fmt"
	"hiv_mind/internal/app"
	"hiv_mind/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("Старт сервиса")

	r := chi.NewRouter()

	serviceProvider := app.NewServiceProvider()
	metricHandler, err := handlers.NewMetricHandler(serviceProvider, "./templates")
	if err != nil {
		return fmt.Errorf("не удалось загрузить шаблоны: %w", err)
	}

	r.Get("/", metricHandler.Root)
	r.Get("/value/{type}/{name}", metricHandler.Value)
	r.Post("/update/{type}/{name}/{value}", metricHandler.Update)

	return http.ListenAndServe(`:8080`, r)
}
