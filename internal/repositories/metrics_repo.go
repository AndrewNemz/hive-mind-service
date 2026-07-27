package repositories

import (
	"fmt"
	"hiv_mind/internal/entities"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (ms *MemStorage) StoreMetric(m entities.Metrics) error {
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
