package repositories

import "hiv_mind/internal/entities"

type IMetricStoragerRepo interface {
	StoreMetric(m *entities.Metrics) error
	StoreMetricSlice(metrics []entities.Metrics) error
	GetAllMetrics() ([]entities.Metrics, error)
	GetMetricByTypeAndName(metric *entities.Metrics) error
	LoadMetricFromFile(filename string) error
	SaveToFile(filename string) error
}
