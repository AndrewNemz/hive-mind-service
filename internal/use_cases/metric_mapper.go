package usecases

import (
	"hiv_mind/internal/entities"
	"runtime"
)

const GaugeType = "gauge"
const CounterType = "counter"

func MapRunTimeMetricToMetricEntitySlice(rm *runtime.MemStats) []entities.Metrics {
	return []entities.Metrics{
		{Type: GaugeType, Name: entities.Alloc, Value: float64(rm.Alloc)},
		{Type: GaugeType, Name: entities.BuckHashSys, Value: float64(rm.BuckHashSys)},
		{Type: GaugeType, Name: entities.Frees, Value: float64(rm.Frees)},
		{Type: GaugeType, Name: entities.GCCPUFraction, Value: rm.GCCPUFraction},
		{Type: GaugeType, Name: entities.GCSys, Value: float64(rm.GCSys)},
		{Type: GaugeType, Name: entities.HeapAlloc, Value: float64(rm.HeapAlloc)},
		{Type: GaugeType, Name: entities.HeapIdle, Value: float64(rm.HeapIdle)},
		{Type: GaugeType, Name: entities.HeapInuse, Value: float64(rm.HeapInuse)},
		{Type: GaugeType, Name: entities.HeapObjects, Value: float64(rm.HeapObjects)},
		{Type: GaugeType, Name: entities.HeapReleased, Value: float64(rm.HeapReleased)},
		{Type: GaugeType, Name: entities.Sys, Value: float64(rm.Sys)},
		{Type: GaugeType, Name: entities.LastGC, Value: float64(rm.LastGC)},
		{Type: GaugeType, Name: entities.Lookups, Value: float64(rm.Lookups)},
		{Type: GaugeType, Name: entities.MCacheInuse, Value: float64(rm.MCacheInuse)},
		{Type: GaugeType, Name: entities.MCacheSys, Value: float64(rm.MCacheSys)},
		{Type: GaugeType, Name: entities.MSpanInuse, Value: float64(rm.MSpanInuse)},
		{Type: GaugeType, Name: entities.MSpanSys, Value: float64(rm.MSpanSys)},
		{Type: GaugeType, Name: entities.Mallocs, Value: float64(rm.Mallocs)},
		{Type: GaugeType, Name: entities.NextGC, Value: float64(rm.NextGC)},
		{Type: GaugeType, Name: entities.NumForcedGC, Value: float64(rm.NumForcedGC)},
		{Type: GaugeType, Name: entities.NumGC, Value: float64(rm.NumGC)},
		{Type: GaugeType, Name: entities.OtherSys, Value: float64(rm.OtherSys)},
		{Type: GaugeType, Name: entities.PauseTotalNs, Value: float64(rm.PauseTotalNs)},
		{Type: GaugeType, Name: entities.StackInuse, Value: float64(rm.StackInuse)},
		{Type: GaugeType, Name: entities.StackSys, Value: float64(rm.StackSys)},
		{Type: GaugeType, Name: entities.Sys, Value: float64(rm.Sys)},
		{Type: GaugeType, Name: entities.TotalAlloc, Value: float64(rm.TotalAlloc)},
	}
}
