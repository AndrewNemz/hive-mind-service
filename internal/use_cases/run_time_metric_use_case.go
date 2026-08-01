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
	MetricRepo   repositories.IMetricStoragerRepo
	MetricSender repositories.IMetricSender
}

func NewRunTimeMetricUseCase(
	metricRepo repositories.IMetricStoragerRepo,
	metricSender repositories.IMetricSender,
) *RunTimeMetricUseCase {
	return &RunTimeMetricUseCase{
		MetricRepo:   metricRepo,
		MetricSender: metricSender,
	}
}

func (rmu *RunTimeMetricUseCase) CollectRunTimeMetric() error {
	rtm := &runtime.MemStats{}
	runtime.ReadMemStats(rtm)

	metrics := MapRunTimeMetricToMetricEntitySlice(rtm)
	if err := rmu.MetricRepo.StoreMetricSlice(metrics); err != nil {
		return err
	}

	return nil
}

func (rmu *RunTimeMetricUseCase) SendRunTimeMetric() error {

	metrics := rmu.MetricRepo.GetAllMetrics()
	if err := rmu.MetricSender.SendMetrics(metrics); err != nil {
		return err
	}

	return nil
}
