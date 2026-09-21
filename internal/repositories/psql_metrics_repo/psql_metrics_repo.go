package psqlmetricsrepo

import "hiv_mind/internal/entities"

type PostgresStorage struct {
}

func NewPostgresStorage() *PostgresStorage {
	return &PostgresStorage{}
}

func (ps *PostgresStorage) StoreMetric(m *entities.Metrics) error {
	return nil
}

func (ps *PostgresStorage) StoreMetricSlice(metrics []entities.Metrics) error {
	return nil
}

func (ps *PostgresStorage) GetAllMetrics() []entities.Metrics {
	return nil
}

func (ps *PostgresStorage) GetMetricByTypeAndName(metric *entities.Metrics) error {
	return nil
}

func (ps *PostgresStorage) LoadMetricFromFile(filename string) error {
	return nil
}

func (ps *PostgresStorage) SaveToFile(filename string) error {
	return nil
}
