package usecases

import (
	"hiv_mind/internal/entities"
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

	var metrics []entities.Metrics
	var statusCode int
	var err error

	err = rmu.BackOff.Retry(
		func() error {
			metrics, err = rmu.MetricRepo.GetAllMetrics()
			return err
		},
		retry.IsRetriableDBError,
	)

	err = rmu.BackOff.Retry(
		func() error {
			statusCode, err = rmu.MetricSender.SendBatchMetrics(metrics)
			return err
		},
		func(err error) bool {
			return retry.IsRetriableHTTPError(statusCode, err)
		},
	)

	return nil
}
