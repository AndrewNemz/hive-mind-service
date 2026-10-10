package usecases

import (
	"hiv_mind/internal/entities"
	"hiv_mind/internal/repositories"
	"hiv_mind/pkg/retry"
)

type IMetricsUseCase interface {
	CollectAndStoreMetric(m *entities.Metrics) error
	GetMetricByTypeAndName(metric *entities.Metrics) error
	GetAllMetrics() ([]entities.Metrics, error)
	UpdateMetricsBatch(metrics []entities.Metrics) error
}

type MetricsUseCase struct {
	repositories repositories.IMetricStoragerRepo
	backoff      *retry.CustomBackOff
}

func NewMetricUseCase(metricsRepo repositories.IMetricStoragerRepo, backoff *retry.CustomBackOff) *MetricsUseCase {
	return &MetricsUseCase{
		repositories: metricsRepo,
		backoff:      backoff,
	}
}

func (mu *MetricsUseCase) CollectAndStoreMetric(m *entities.Metrics) error {
	err := mu.backoff.Retry(
		func() error {
			return mu.repositories.StoreMetric(m)
		},
		retry.IsRetriableDBError,
	)

	return err
}

func (mu *MetricsUseCase) GetMetricByTypeAndName(metric *entities.Metrics) error {
	err := mu.backoff.Retry(
		func() error {
			return mu.repositories.GetMetricByTypeAndName(metric)
		},
		retry.IsRetriableDBError,
	)

	return err
}

func (mu *MetricsUseCase) GetAllMetrics() []entities.Metrics {
	var metrics []entities.Metrics
	var err error

	err = mu.backoff.Retry(
		func() error {
			metrics, err = mu.repositories.GetAllMetrics()
			return err
		},
		retry.IsRetriableDBError,
	)

	return metrics
}

func (mu *MetricsUseCase) UpdateMetricsBatch(metrics []entities.Metrics) error {
	if err := mu.repositories.StoreMetricSlice(metrics); err != nil {
		return err
	}
	return nil
}
