package usecases

import (
	"hiv_mind/internal/entities"
	"hiv_mind/internal/repositories"
)

type IMetricsUseCase interface {
	CollectAndStoreMetric(m entities.Metrics) error
	GetMetricByTypeAndName(mType, mName string) (entities.Metrics, error)
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

func (mu *MetricsUseCase) CollectAndStoreMetric(m entities.Metrics) error {
	if err := mu.repositories.StoreMetric(m); err != nil {
		return err
	}

	return nil
}

func (mu *MetricsUseCase) GetMetricByTypeAndName(mType, mName string) (entities.Metrics, error) {
	metric, err := mu.repositories.GetMetricByTypeAndName(mType, mName)
	if err != nil {
		return entities.Metrics{}, err
	}
	return metric, nil
}

func (mu *MetricsUseCase) GetAllMetrics() []entities.Metrics {
	metrics := mu.repositories.GetAllMetrics()
	return metrics
}
