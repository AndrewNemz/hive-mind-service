package usecases

import (
	"hiv_mind/internal/entities"
	"runtime"
)

func MapRunTimeMetricToMetricEntitySlice(rm *runtime.MemStats) []entities.Metrics {
	return []entities.Metrics{
		{Type: entities.GaugeType, Name: entities.Alloc, Value: float64(rm.Alloc)},
		{Type: entities.GaugeType, Name: entities.BuckHashSys, Value: float64(rm.BuckHashSys)},
		{Type: entities.GaugeType, Name: entities.Frees, Value: float64(rm.Frees)},
		{Type: entities.GaugeType, Name: entities.GCCPUFraction, Value: rm.GCCPUFraction},
		{Type: entities.GaugeType, Name: entities.GCSys, Value: float64(rm.GCSys)},
		{Type: entities.GaugeType, Name: entities.HeapAlloc, Value: float64(rm.HeapAlloc)},
		{Type: entities.GaugeType, Name: entities.HeapIdle, Value: float64(rm.HeapIdle)},
		{Type: entities.GaugeType, Name: entities.HeapInuse, Value: float64(rm.HeapInuse)},
		{Type: entities.GaugeType, Name: entities.HeapObjects, Value: float64(rm.HeapObjects)},
		{Type: entities.GaugeType, Name: entities.HeapReleased, Value: float64(rm.HeapReleased)},
		{Type: entities.GaugeType, Name: entities.Sys, Value: float64(rm.Sys)},
		{Type: entities.GaugeType, Name: entities.LastGC, Value: float64(rm.LastGC)},
		{Type: entities.GaugeType, Name: entities.Lookups, Value: float64(rm.Lookups)},
		{Type: entities.GaugeType, Name: entities.MCacheInuse, Value: float64(rm.MCacheInuse)},
		{Type: entities.GaugeType, Name: entities.MCacheSys, Value: float64(rm.MCacheSys)},
		{Type: entities.GaugeType, Name: entities.MSpanInuse, Value: float64(rm.MSpanInuse)},
		{Type: entities.GaugeType, Name: entities.MSpanSys, Value: float64(rm.MSpanSys)},
		{Type: entities.GaugeType, Name: entities.Mallocs, Value: float64(rm.Mallocs)},
		{Type: entities.GaugeType, Name: entities.NextGC, Value: float64(rm.NextGC)},
		{Type: entities.GaugeType, Name: entities.NumForcedGC, Value: float64(rm.NumForcedGC)},
		{Type: entities.GaugeType, Name: entities.NumGC, Value: float64(rm.NumGC)},
		{Type: entities.GaugeType, Name: entities.OtherSys, Value: float64(rm.OtherSys)},
		{Type: entities.GaugeType, Name: entities.PauseTotalNs, Value: float64(rm.PauseTotalNs)},
		{Type: entities.GaugeType, Name: entities.StackInuse, Value: float64(rm.StackInuse)},
		{Type: entities.GaugeType, Name: entities.StackSys, Value: float64(rm.StackSys)},
		{Type: entities.GaugeType, Name: entities.Sys, Value: float64(rm.Sys)},
		{Type: entities.GaugeType, Name: entities.TotalAlloc, Value: float64(rm.TotalAlloc)},
	}
}
