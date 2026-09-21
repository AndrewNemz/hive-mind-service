package psqlmetricsrepo

import (
	"context"
	"fmt"
	"hiv_mind/internal/entities"
	"hiv_mind/pkg/logger"
	"hiv_mind/pkg/postgresql"
)

type PostgresStorage struct {
	db *postgresql.PostgreSQL
}

func NewPostgresStorage(db *postgresql.PostgreSQL) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (ps *PostgresStorage) StoreMetric(m *entities.Metrics) error {

	lg := logger.Get()
	if ps == nil || ps.db == nil || ps.db.DB == nil {
		lg.Error("хранилище PostgreSQL не инициализировано (db is nil). Запустите сервер с флагом -d")
		return fmt.Errorf("хранилище PostgreSQL не инициализировано (db is nil). Запустите сервер с флагом -d")
	}

	query := `
		INSERT INTO metrics (id, type, value, delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id, type) DO UPDATE
		SET 
			value = CASE WHEN EXCLUDED.value IS NOT NULL THEN EXCLUDED.value ELSE metrics.value END,
			delta = CASE WHEN EXCLUDED.delta IS NOT NULL THEN metrics.delta + EXCLUDED.delta ELSE metrics.delta END
	`

	// TODO вынести контекст в хандлеру
	_, err := ps.db.DB.Exec(
		context.Background(),
		query,
		m.ID,
		m.MType,
		m.Value,
		m.Delta,
	)

	if err != nil {
		return err
	}

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
