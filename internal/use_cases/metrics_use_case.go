package usecases

import (
	"hiv_mind/internal/entities"
	"hiv_mind/internal/repositories"
)

type IMetricsUseCase interface {
	CollectAndStoreMetric(m entities.Metrics) error
}

type MetricsUseCase struct {
	repositories repositories.IMetricStoragerRepo
}

func NewMetricUseCase(metricsRepo repositories.IMetricStoragerRepo) *MetricsUseCase {
	return &MetricsUseCase{
		repositories: metricsRepo,
	}
}

func (mu *MetricsUseCase) CollectAndStoreMetric(m entities.Metrics) error {
	if err := mu.repositories.StoreMetric(m); err != nil {
		return err
	}

	return nil
}
