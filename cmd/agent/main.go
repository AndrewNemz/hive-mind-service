package main

import (
	"flag"
	"fmt"
	"hiv_mind/internal/app"
	"hiv_mind/internal/handlers"
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

	if err := hundler.HandleRunTimeMetric(); err != nil {
		return err
	}

	return nil
}
