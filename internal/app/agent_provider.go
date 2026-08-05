package app

import (
	"hiv_mind/internal/repositories"
	usecases "hiv_mind/internal/use_cases"
)

type AgentProvider struct {
	Storage              repositories.IMetricStoragerRepo
	RunTimeMetricUseCase usecases.IRunTimeMetricUseCase
	MetricSender         repositories.IMetricSender

	Adresss        string
	ReportInterval int
	PollInterval   int
}

func NewAgentProvider(adresss string, reportInterval, pollInterval int) *AgentProvider {
	storage := repositories.NewMemStorage()
	sender := repositories.NewMetricSender(adresss)
	metricUseCase := usecases.NewRunTimeMetricUseCase(storage, sender)
	return &AgentProvider{
		Storage:              storage,
		RunTimeMetricUseCase: metricUseCase,
		Adresss:              adresss,
		ReportInterval:       reportInterval,
		PollInterval:         pollInterval,
	}
}
