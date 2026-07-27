package repositories

import "hiv_mind/internal/entities"

type IMetricStoragerRepo interface {
	StoreMetric(m entities.Metrics) error
}
