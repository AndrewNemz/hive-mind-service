package main

import (
	"flag"
	"fmt"
	"hiv_mind/internal/app"
	"hiv_mind/internal/handlers"
	"hiv_mind/pkg/logger"

	"go.uber.org/zap"
)

var (
	reportInterval *int
	pollInterval   *int
	adresss        *string
)

func init() {
	reportInterval = flag.Int("r", 10, "частота отправки метрик на сервер (по умолчанию 10 секунд)")
	pollInterval = flag.Int("p", 2, "частоту опроса метрик из пакета runtime (по умолчанию 2 секунды)")
	adresss = flag.String("a", "localhost:8080", "отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080)")
}

func main() {
	fmt.Println("Старт агента")

	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()
	agentProvider := app.NewAgentProvider(*adresss, *reportInterval, *pollInterval)
	hundler := handlers.NewAgentHandler(agentProvider)

	if err := logger.Initialize("info"); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	lg := logger.Get()
	defer lg.Sync()

	if err := hundler.HandleRunTimeMetric(); err != nil {
		lg.Info("Failed in Handle Request", zap.Error(err))
		return err
	}

	return nil
}
