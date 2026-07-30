package app

import (
	"hiv_mind/internal/repositories"
	usecases "hiv_mind/internal/use_cases"
)

type AgentProvider struct {
	Storage              repositories.IMetricStoragerRepo
	RunTimeMetricUseCase usecases.IRunTimeMetricUseCase
	MetricSender         repositories.IMetricSender
}

func NewAgentProvider() *AgentProvider {
	storage := repositories.NewMemStorage()
	sender := repositories.NewMetricSender()
	metricUseCase := usecases.NewRunTimeMetricUseCase(storage, sender)
	return &AgentProvider{
		Storage:              storage,
		RunTimeMetricUseCase: metricUseCase,
	}
}
