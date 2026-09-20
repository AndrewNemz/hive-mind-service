package main

import (
	"flag"
	"fmt"
	"hiv_mind/internal/app"
	"hiv_mind/internal/handlers"
	"hiv_mind/pkg/logger"
	"hiv_mind/pkg/middleware"
	"net/http"
	"os"

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

	if err := logger.Initialize("info"); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	lg := logger.Get()
	defer lg.Sync()

	// middlewares
	r.Use(middleware.RequestInfo)
	r.Use(middleware.ResponseInfo)

	// Routes
	r.Get("/", metricHandler.Root)
	r.Get("/value/{type}/{name}", metricHandler.Value)
	r.Post("/update/{type}/{name}/{value}", metricHandler.Update)

	if envAdress := os.Getenv("ADDRESS"); envAdress != "" {
		adresss = envAdress
	}
	return http.ListenAndServe(adresss, r)
}
