package usecases

import (
	"hiv_mind/internal/entities"
	"runtime"
)

func MapRunTimeMetricToMetricEntitySlice(rm *runtime.MemStats) []entities.Metrics {
	type metricMapping struct {
		id    string
		value float64
	}

	mappings := []metricMapping{
		{entities.Alloc, float64(rm.Alloc)},
		{entities.BuckHashSys, float64(rm.BuckHashSys)},
		{entities.Frees, float64(rm.Frees)},
		{entities.GCCPUFraction, rm.GCCPUFraction},
		{entities.GCSys, float64(rm.GCSys)},
		{entities.HeapAlloc, float64(rm.HeapAlloc)},
		{entities.HeapIdle, float64(rm.HeapIdle)},
		{entities.HeapInuse, float64(rm.HeapInuse)},
		{entities.HeapObjects, float64(rm.HeapObjects)},
		{entities.HeapReleased, float64(rm.HeapReleased)},
		{entities.HeapSys, float64(rm.HeapSys)},
		{entities.Sys, float64(rm.Sys)},
		{entities.LastGC, float64(rm.LastGC)},
		{entities.Lookups, float64(rm.Lookups)},
		{entities.MCacheInuse, float64(rm.MCacheInuse)},
		{entities.MCacheSys, float64(rm.MCacheSys)},
		{entities.MSpanInuse, float64(rm.MSpanInuse)},
		{entities.MSpanSys, float64(rm.MSpanSys)},
		{entities.Mallocs, float64(rm.Mallocs)},
		{entities.NextGC, float64(rm.NextGC)},
		{entities.NumForcedGC, float64(rm.NumForcedGC)},
		{entities.NumGC, float64(rm.NumGC)},
		{entities.OtherSys, float64(rm.OtherSys)},
		{entities.PauseTotalNs, float64(rm.PauseTotalNs)},
		{entities.StackInuse, float64(rm.StackInuse)},
		{entities.StackSys, float64(rm.StackSys)},
		{entities.Sys, float64(rm.Sys)},
		{entities.TotalAlloc, float64(rm.TotalAlloc)},
	}

	result := make([]entities.Metrics, 0, len(mappings))
	for _, m := range mappings {
		result = append(result, entities.Metrics{
			ID:    m.id,
			MType: entities.GaugeType,
			Value: &m.value,
		})
	}

	return result
}
