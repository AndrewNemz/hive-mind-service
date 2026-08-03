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

func (ms *MemStorage) StoreMetric(m entities.Metrics) error {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()

	switch m.Type {
	case "gauge":
		ms.Gauge[m.Name] = m.Value
	case "counter":
		ms.Counter[m.Name] += int64(m.Value)
	default:
		return fmt.Errorf("В репозиторий передан неожиданный формат метрики!")
	}

	return nil
}

func (ms *MemStorage) StoreMetricSlice(metrics []entities.Metrics) error {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()

	for _, m := range metrics {
		switch m.Type {
		case "gauge":
			ms.Gauge[m.Name] = m.Value
		case "counter":
			ms.Counter[m.Name] += int64(m.Value)
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
		metrics = append(metrics, entities.Metrics{Type: "gauge", Name: metricName, Value: gm})
	}

	for metricName, cm := range ms.Counter {
		metrics = append(metrics, entities.Metrics{Type: "counter", Name: metricName, Value: float64(cm)})
	}

	return metrics
}

func (ms *MemStorage) GetMetricByTypeAndName(mType, mName string) (entities.Metrics, error) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()

	var mValue float64
	var valueI int64
	var ok bool

	if mType == entities.GaugeType {
		mValue, ok = ms.Gauge[mName]
	} else {
		valueI, ok = ms.Counter[mName]
		mValue = float64(valueI)
	}
	if !ok {
		return entities.Metrics{}, repoerrors.ErrNotFoundMetric
	}

	return entities.Metrics{Name: mName, Type: mType, Value: mValue}, nil
}
