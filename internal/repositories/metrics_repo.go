package repositories

import (
	"encoding/json"
	"fmt"
	"hiv_mind/internal/entities"
	repoerrors "hiv_mind/internal/errors"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
	Mutex   sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
		Mutex:   sync.Mutex{},
	}
}

func (ms *MemStorage) StoreMetric(m *entities.Metrics) error {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()

	switch m.MType {
	case "gauge":
		ms.Gauge[m.ID] = *m.Value
	case "counter":
		ms.Counter[m.ID] += int64(*m.Delta)
	default:
		return fmt.Errorf("В репозиторий передан неожиданный формат метрики!")
	}

	return nil
}

func (ms *MemStorage) StoreMetricSlice(metrics []entities.Metrics) error {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()

	for _, m := range metrics {
		switch m.MType {
		case "gauge":
			ms.Gauge[m.ID] = *m.Value
		case "counter":
			ms.Counter[m.ID] += int64(*m.Delta)
		default:
			return fmt.Errorf("В репозиторий передан неожиданный формат метрики!")
		}
	}

	// TODO: В продакшене логику инкремента PollCount и генерации RandomValue
	// следует вынести в UseCase, чтобы репозиторий оставался "глупым" хранилищем.
	// Сейчас это допустимо, так как UseCase не должен знать о внутреннем состоянии хранилища.
	ms.Counter[entities.PollCount]++
	ms.Gauge[entities.RandomValue] = rand.Float64()

	return nil
}

func (ms *MemStorage) GetAllMetrics() ([]entities.Metrics, error) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()

	var metrics []entities.Metrics
	for metricName, gm := range ms.Gauge {
		val := gm
		metrics = append(metrics, entities.Metrics{MType: "gauge", ID: metricName, Value: &val})
	}

	for metricName, cm := range ms.Counter {
		val := cm
		metrics = append(metrics, entities.Metrics{MType: "counter", ID: metricName, Delta: &val})
	}

	return metrics, nil
}

func (ms *MemStorage) GetMetricByTypeAndName(metric *entities.Metrics) error {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()

	var mValue float64
	var valueI int64
	var ok bool

	if metric.MType == entities.GaugeType {
		mValue, ok = ms.Gauge[metric.ID]
		if !ok {
			return repoerrors.ErrNotFoundMetric
		}
		metric.Value = &mValue
	} else {
		valueI, ok = ms.Counter[metric.ID]
		if !ok {
			return repoerrors.ErrNotFoundMetric
		}
		metric.Delta = &valueI
	}

	return nil
}

func (ms *MemStorage) LoadMetricFromFile(filename string) error {
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

	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()

	for k, v := range fileData.Gauge {
		ms.Gauge[k] = v
	}

	for k, v := range fileData.Counter {
		ms.Counter[k] = v
	}

	return nil
}

func (ms *MemStorage) SaveToFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("Файл не указан!")
	}

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	ms.Mutex.Lock()

	data := struct {
		Gauge   map[string]float64 `json:"gauge"`
		Counter map[string]int64   `json:"counter"`
	}{
		Gauge:   make(map[string]float64, len(ms.Gauge)),
		Counter: make(map[string]int64, len(ms.Counter)),
	}

	for k, v := range ms.Gauge {
		data.Gauge[k] = v
	}
	for k, v := range ms.Counter {
		data.Counter[k] = v
	}

	ms.Mutex.Unlock()

	tempFile := filename + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	// defer закрытия файла. Если будет ошибка, временный файл останется,
	// он перезапишется при следующей попытке.
	defer func() {
		file.Close()
		// Если произошла ошибка, удаляем битый временный файл
		if err != nil {
			os.Remove(tempFile)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode metrics to json: %w", err)
	}

	if err = file.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err = os.Rename(tempFile, filename); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}
