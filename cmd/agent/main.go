package main

import (
	"fmt"
	"hiv_mind/internal/app"
	"hiv_mind/internal/handlers"
)

func main() {
	fmt.Println("Старт агента")

	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	agentProvider := app.NewAgentProvider()
	hundler := handlers.NewAgentHandler(agentProvider)

	if err := hundler.HandleRunTimeMetric(); err != nil {
		return err
	}

	return nil
}
