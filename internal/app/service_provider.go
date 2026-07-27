package app

import (
	"hiv_mind/internal/repositories"
	usecases "hiv_mind/internal/use_cases"
)

type ServiceProvider struct {
	Storage       repositories.IMetricStoragerRepo
	MetricUseCase usecases.IMetricsUseCase
}

func NewServiceProvider() *ServiceProvider {
	storage := repositories.NewMemStorage()
	metricUseCase := usecases.NewMetricUseCase(storage)
	return &ServiceProvider{
		Storage:       storage,
		MetricUseCase: metricUseCase,
	}
}
