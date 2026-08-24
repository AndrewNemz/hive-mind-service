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
	"go.uber.org/zap"
)

var adress string

func init() {
	flag.StringVar(&adress, "a", "localhost:8080", "отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080)")
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize("info"); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	lg := logger.Get()
	defer lg.Sync()

	flag.Parse()
	lg.Sugar().Info("Старт сервиса", zap.String("adress", adress))

	r := chi.NewRouter()
	serviceProvider := app.NewServiceProvider()
	metricHandler, err := handlers.NewMetricHandler(serviceProvider, "./templates")
	if err != nil {
		lg.Error("не удалось загрузить шаблоны", zap.Error(err))
		return err
	}

	// middlewares
	r.Use(middleware.Logging)

	// Routes
	r.Get("/", metricHandler.Root)
	r.Post("/value/", metricHandler.Value)
	r.Post("/update/", metricHandler.Update)

	if envAdress := os.Getenv("ADDRESS"); envAdress != "" {
		adress = envAdress
	}
	return http.ListenAndServe(adress, r)
}
