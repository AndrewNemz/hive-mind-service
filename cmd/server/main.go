package main

import (
	"fmt"
	"hiv_mind/internal/app"
	"hiv_mind/internal/handlers"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("Старт сервиса")
	serviceProvider := app.NewServiceProvider()
	metricHandler := handlers.NewMetricHandler(serviceProvider)

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", metricHandler.Update)

	return http.ListenAndServe(`:8080`, mux)
}
