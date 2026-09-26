package psqlmetricsrepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hiv_mind/internal/entities"
	repoerrors "hiv_mind/internal/errors"
	"hiv_mind/pkg/logger"
	"hiv_mind/pkg/postgresql"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
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
	lg := logger.Get()
	if ps == nil || ps.db == nil || ps.db.DB == nil {
		lg.Error("хранилище PostgreSQL не инициализировано (db is nil). Запустите сервер с флагом -d")
		return fmt.Errorf("хранилище PostgreSQL не инициализировано (db is nil). Запустите сервер с флагом -d")
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO metrics (id, type, value, delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id, type) DO UPDATE
		SET 
			value = CASE WHEN EXCLUDED.value IS NOT NULL THEN EXCLUDED.value ELSE metrics.value END,
			delta = CASE WHEN EXCLUDED.delta IS NOT NULL THEN metrics.delta + EXCLUDED.delta ELSE metrics.delta END
	`
	for _, m := range metrics {
		batch.Queue(query, m.ID, m.MType, m.Value, m.Delta)
	}

	batchR := ps.db.DB.SendBatch(context.Background(), batch)
	defer batchR.Close()

	for i := 0; i < len(metrics); i++ {
		_, err := batchR.Exec()
		if err != nil {
			lg.Error(
				"Ошибка сохранения метрики из батча",
				zap.Int("index", i),
				zap.Error(err),
			)
		}
	}

	return nil
}

func (ps *PostgresStorage) GetAllMetrics() ([]entities.Metrics, error) {
	lg := logger.Get()
	if ps == nil || ps.db == nil || ps.db.DB == nil {
		lg.Error("хранилище PostgreSQL не инициализировано (db is nil). Запустите сервер с флагом -d")
		return nil, fmt.Errorf("хранилище PostgreSQL не инициализировано (db is nil). Запустите сервер с флагом -d")
	}

	metrics := []entities.Metrics{}
	query := `SELECT id, type, value, delta FROM metrics`

	rows, err := ps.db.DB.Query(context.Background(), query)
	if err != nil {
		lg.Error("GetAllMetrics() - Ошибка при получении метрик!", zap.Error(err))
		return nil, fmt.Errorf("Ошибка при получении метрик!")
	}
	defer rows.Close()

	for rows.Next() {
		var m entities.Metrics
		var valuePtr *float64
		var deltaPtr *int64

		err := rows.Scan(&m.ID, &m.MType, valuePtr, deltaPtr)
		if err != nil {
			lg.Error("GetAllMetrics() - Ошибка при чтении колонок rows!", zap.Error(err))
			return nil, fmt.Errorf("Ошибка при получении метрик!")
		}

		m.Value = valuePtr
		m.Delta = deltaPtr
		metrics = append(metrics, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (ps *PostgresStorage) GetMetricByTypeAndName(metric *entities.Metrics) error {
	lg := logger.Get()
	if ps == nil || ps.db == nil || ps.db.DB == nil {
		lg.Error("хранилище PostgreSQL не инициализировано (db is nil). Запустите сервер с флагом -d")
		return fmt.Errorf("хранилище PostgreSQL не инициализировано (db is nil). Запустите сервер с флагом -d")
	}

	query := `SELECT value, delta FROM metrics WHERE id=$1 AND type=$2`
	var valuePtr *float64
	var deltaPtr *int64

	if err := ps.db.DB.QueryRow(
		context.Background(), query, metric.ID, metric.MType,
	).Scan(valuePtr, deltaPtr); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repoerrors.ErrNotFoundMetric
		}
		return fmt.Errorf("GetMetricByTypeAndName - Ошибка при получении метрики")
	}

	metric.Delta = deltaPtr
	metric.Value = valuePtr
	return nil
}

func (ps *PostgresStorage) LoadMetricFromFile(filename string) error {
	if filename == "" {
		return nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read metrics file: %w", err)
	}

	var fileData struct {
		Gauge   map[string]float64 `json:"gauge"`
		Counter map[string]int64   `json:"counter"`
	}

	if err := json.Unmarshal(data, &fileData); err != nil {
		return fmt.Errorf("failed to unmarshal metrics file: %w", err)
	}

	var metrics []entities.Metrics

	for k, v := range fileData.Gauge {
		val := v
		metrics = append(metrics, entities.Metrics{
			ID:    k,
			MType: entities.GaugeType,
			Value: &val,
		})
	}

	for k, v := range fileData.Counter {
		val := v
		metrics = append(metrics, entities.Metrics{
			ID:    k,
			MType: entities.CounterType,
			Delta: &val,
		})
	}

	return ps.StoreMetricSlice(metrics)
}

func (ps *PostgresStorage) SaveToFile(filename string) error {
	if filename == "" {
		return nil
	}

	// 1 Получаем актуальные данные из БД
	metrics, err := ps.GetAllMetrics()
	if err != nil {
		return fmt.Errorf("failed to get metrics for saving: %w", err)
	}

	// 2 Создаем директорию, если её нет
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// 3 Конвертируем срез в структуру для JSON (аналогично MemStorage)
	data := struct {
		Gauge   map[string]float64 `json:"gauge"`
		Counter map[string]int64   `json:"counter"`
	}{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}

	for _, m := range metrics {
		if m.MType == entities.GaugeType && m.Value != nil {
			data.Gauge[m.ID] = *m.Value
		} else if m.MType == entities.CounterType && m.Delta != nil {
			data.Counter[m.ID] = *m.Delta
		}
	}

	// 4 Запись в файл
	tempFile := filename + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	var writeErr error
	defer func() {
		file.Close()
		if writeErr != nil {
			os.Remove(tempFile)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(data); err != nil {
		writeErr = fmt.Errorf("failed to encode metrics to json: %w", err)
		return writeErr
	}

	if err = file.Sync(); err != nil {
		writeErr = fmt.Errorf("failed to sync file: %w", err)
		return writeErr
	}

	if err = file.Close(); err != nil {
		writeErr = fmt.Errorf("failed to close temp file: %w", err)
		return writeErr
	}

	if err = os.Rename(tempFile, filename); err != nil {
		writeErr = fmt.Errorf("failed to rename temp file: %w", err)
		return writeErr
	}

	return nil
}
