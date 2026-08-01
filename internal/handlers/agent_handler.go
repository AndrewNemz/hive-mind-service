package handlers

import (
	"context"
	"hiv_mind/internal/app"
	"sync"
	"time"
)

const SleepTime = 2
const reportInterval = 10
const pollInterval = 2

type AgentHandler struct {
	ap *app.AgentProvider
}

func NewAgentHandler(ap *app.AgentProvider) *AgentHandler {
	return &AgentHandler{
		ap: ap,
	}
}

func (ah *AgentHandler) HandleRunTimeMetric() error {

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		ticker := time.NewTicker(pollInterval * time.Second)
		defer ticker.Stop()

		for {
			if err := ah.ap.RunTimeMetricUseCase.CollectRunTimeMetric(); err != nil {
				cancelFunc()
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()

	go func() {
		defer wg.Done()

		ticker := time.NewTicker(reportInterval * time.Second)
		defer ticker.Stop()

		for {
			if err := ah.ap.RunTimeMetricUseCase.SendRunTimeMetric(); err != nil {
				cancelFunc()
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()

	wg.Wait()
	return nil
}
