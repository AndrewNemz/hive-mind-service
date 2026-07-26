package main

import (
	"fmt"
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

	storage := handlers.NewMemStorage()
	return http.ListenAndServe(`:8080`, http.StripPrefix("/update/", storage))
}
