package usecases

import (
	"hiv_mind/internal/entities"
	"hiv_mind/internal/repositories"
)

type IMetricsUseCase interface {
	CollectAndStoreMetric(m *entities.Metrics) error
	GetMetricByTypeAndName(metric *entities.Metrics) error
	GetAllMetrics() []entities.Metrics
}

type MetricsUseCase struct {
	repositories repositories.IMetricStoragerRepo
}

func NewMetricUseCase(metricsRepo repositories.IMetricStoragerRepo) *MetricsUseCase {
	return &MetricsUseCase{
		repositories: metricsRepo,
	}
}

func (mu *MetricsUseCase) CollectAndStoreMetric(m *entities.Metrics) error {
	if err := mu.repositories.StoreMetric(m); err != nil {
		return err
	}

	return nil
}

func (mu *MetricsUseCase) GetMetricByTypeAndName(metric *entities.Metrics) error {
	err := mu.repositories.GetMetricByTypeAndName(metric)
	if err != nil {
		return err
	}
	return nil
}

func (mu *MetricsUseCase) GetAllMetrics() []entities.Metrics {
	metrics := mu.repositories.GetAllMetrics()
	return metrics
}
