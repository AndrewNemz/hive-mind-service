package repositories

import "hiv_mind/internal/entities"

type IMetricSender interface {
	SendMetrics(metrics []entities.Metrics) error
}
