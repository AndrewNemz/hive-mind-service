package app

import (
	"hiv_mind/internal/repositories"
	usecases "hiv_mind/internal/use_cases"
)

type AgentProvider struct {
	Storage              repositories.IMetricStoragerRepo
	RunTimeMetricUseCase usecases.IRunTimeMetricUseCase
}

func NewAgentProvider() *AgentProvider {
	storage := repositories.NewMemStorage()
	metricUseCase := usecases.NewRunTimeMetricUseCase(storage)
	return &AgentProvider{
		Storage:              storage,
		RunTimeMetricUseCase: metricUseCase,
	}
}
