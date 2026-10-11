package app

import (
	"hiv_mind/internal/repositories"
	usecases "hiv_mind/internal/use_cases"
	"hiv_mind/pkg/retry"
)

type AgentProvider struct {
	Storage              repositories.IMetricStoragerRepo
	RunTimeMetricUseCase usecases.IRunTimeMetricUseCase
	MetricSender         repositories.IMetricSender

	Adresss        string
	ReportInterval int
	PollInterval   int
}

func NewAgentProvider(adresss string, reportInterval, pollInterval int, hashKey string) *AgentProvider {
	storage := repositories.NewMemStorage()
	sender := repositories.NewMetricSender(adresss, hashKey)
	backoff := retry.NewCustomBackOff()
	metricUseCase := usecases.NewRunTimeMetricUseCase(storage, sender, backoff)
	return &AgentProvider{
		Storage:              storage,
		RunTimeMetricUseCase: metricUseCase,
		Adresss:              adresss,
		ReportInterval:       reportInterval,
		PollInterval:         pollInterval,
	}
}
