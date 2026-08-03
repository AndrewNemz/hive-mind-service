package repositories

import "hiv_mind/internal/entities"

type IMetricStoragerRepo interface {
	StoreMetric(m entities.Metrics) error
	StoreMetricSlice(metrics []entities.Metrics) error
	GetAllMetrics() []entities.Metrics
	GetMetricByTypeAndName(mType, mName string) (entities.Metrics, error)
}
