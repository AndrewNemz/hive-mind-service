package usecases

import (
	"hiv_mind/internal/repositories"
	"runtime"
)

type IRunTimeMetricUseCase interface {
	CollectRunTimeMetric() error
	SendRunTimeMetric() error
}

type RunTimeMetricUseCase struct {
	RunTimeMetricRepo repositories.IMetricStoragerRepo
}

func NewRunTimeMetricUseCase(metricRepo repositories.IMetricStoragerRepo) *RunTimeMetricUseCase {
	return &RunTimeMetricUseCase{
		RunTimeMetricRepo: metricRepo,
	}
}

func (rmu *RunTimeMetricUseCase) CollectRunTimeMetric() error {
	rtm := &runtime.MemStats{}
	runtime.ReadMemStats(rtm)

	metrics := MapRunTimeMetricToMetricEntitySlice(rtm)
	if err := rmu.RunTimeMetricRepo.StoreMetricSlice(metrics); err != nil {
		return err
	}

	return nil
}

func (rmu *RunTimeMetricUseCase) SendRunTimeMetric() error {
	return nil
}
