package usecases

import (
	"hiv_mind/internal/repositories"
	"hiv_mind/pkg/retry"
	"runtime"
)

type IRunTimeMetricUseCase interface {
	CollectRunTimeMetric() error
	SendRunTimeMetric() error
}

type RunTimeMetricUseCase struct {
	MetricRepo   repositories.IMetricStoragerRepo
	MetricSender repositories.IMetricSender
	BackOff      *retry.CustomBackOff
}

func NewRunTimeMetricUseCase(
	metricRepo repositories.IMetricStoragerRepo,
	metricSender repositories.IMetricSender,
	backoff *retry.CustomBackOff,
) *RunTimeMetricUseCase {
	return &RunTimeMetricUseCase{
		MetricRepo:   metricRepo,
		MetricSender: metricSender,
		BackOff:      backoff,
	}
}

func (rmu *RunTimeMetricUseCase) CollectRunTimeMetric() error {
	rtm := &runtime.MemStats{}
	runtime.ReadMemStats(rtm)

	metrics := MapRunTimeMetricToMetricEntitySlice(rtm)

	err := rmu.BackOff.Retry(
		func() error {
			return rmu.MetricRepo.StoreMetricSlice(metrics)
		},
		retry.IsRetriableDBError,
	)

	return err
}

func (rmu *RunTimeMetricUseCase) SendRunTimeMetric() error {

	metrics, err := rmu.MetricRepo.GetAllMetrics()
	if err != nil {
		return err
	}
	if err := rmu.MetricSender.SendBatchMetrics(metrics); err != nil {
		return err
	}

	return nil
}
