package handlers

import (
	"hiv_mind/internal/app"
)

const SleepTime = 2

type AgentHandler struct {
	ap *app.AgentProvider
}

func NewAgentHandler(ap *app.AgentProvider) *AgentHandler {
	return &AgentHandler{
		ap: ap,
	}
}

func (ah *AgentHandler) HandleRunTimeMetric() error {

	if err := ah.ap.RunTimeMetricUseCase.CollectRunTimeMetric(); err != nil {
		return err
	}

	return nil
}
