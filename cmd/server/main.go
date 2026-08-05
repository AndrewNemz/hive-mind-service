package main

import (
	"flag"
	"fmt"
	"hiv_mind/internal/app"
	"hiv_mind/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var adresss string

func init() {
	flag.StringVar(&adresss, "a", "localhost:8080", "отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080)")
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("Старт сервиса")
	flag.Parse()

	r := chi.NewRouter()

	serviceProvider := app.NewServiceProvider()
	metricHandler, err := handlers.NewMetricHandler(serviceProvider, "./templates")
	if err != nil {
		return fmt.Errorf("не удалось загрузить шаблоны: %w", err)
	}

	r.Get("/", metricHandler.Root)
	r.Get("/value/{type}/{name}", metricHandler.Value)
	r.Post("/update/{type}/{name}/{value}", metricHandler.Update)

	return http.ListenAndServe(adresss, r)
}
