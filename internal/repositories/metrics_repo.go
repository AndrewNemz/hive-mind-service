package repositories

import (
	"fmt"
	"hiv_mind/internal/entities"
	repoerrors "hiv_mind/internal/errors"
	"math/rand"
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

func (ms *MemStorage) GetAllMetrics() []entities.Metrics {
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

	return metrics
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
