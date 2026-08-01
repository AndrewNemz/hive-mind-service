package usecases

import (
	"fmt"
	"hiv_mind/internal/entities"
	mocksrepo "hiv_mind/internal/repositories/mocks_repo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunTimeMetricUseCaseMethodCollectRunTimeMetric(t *testing.T) {
	data := []struct {
		tName   string
		metrics []entities.Metrics
		err     error
	}{
		{
			tName: "Test_CollectRunTimeMetric_When_OK",
			metrics: []entities.Metrics{
				{Type: entities.CounterType, Name: "Alloc", Value: 10},
				{Type: entities.GaugeType, Name: "Alloc", Value: 12.0},
			},
			err: nil,
		},
		{
			tName: "Test_CollectRunTimeMetric_Whith_Error",
			err:   fmt.Errorf("Repo Error!"),
		},
	}

	for _, d := range data {
		t.Run(d.tName, func(t *testing.T) {
			metricRepo := &mocksrepo.MockStorageRepo{
				HasCalled: true,
				Metrics:   d.metrics,
				Err:       d.err,
			}
			metricSender := &mocksrepo.MockSenderRepo{
				HasCalled: true,
				Metrics:   d.metrics,
				Err:       d.err,
			}
			usecase := NewRunTimeMetricUseCase(metricRepo, metricSender)

			err := usecase.CollectRunTimeMetric()
			if err != nil {
				assert.EqualError(t, d.err, err.Error())
			}
			assert.Equal(t, d.err, err)
		})
	}
}

func TestSendRunTimeMetric(t *testing.T) {
	data := []struct {
		tName   string
		metrics []entities.Metrics
		err     error
	}{
		{
			tName: "Test_SendRunTimeMetric_When_OK",
			metrics: []entities.Metrics{
				{Type: entities.CounterType, Name: "Alloc", Value: 10},
				{Type: entities.GaugeType, Name: "Alloc", Value: 12.0},
			},
			err: nil,
		},
		{
			tName: "Test_SendRunTimeMetric_Whith_Error",
			err:   fmt.Errorf("Repo Error!"),
		},
	}

	for _, d := range data {
		t.Run(d.tName, func(t *testing.T) {
			metricRepo := &mocksrepo.MockStorageRepo{
				Metrics: d.metrics,
				Err:     d.err,
			}
			metricSender := &mocksrepo.MockSenderRepo{
				Metrics: d.metrics,
				Err:     d.err,
			}
			usecase := NewRunTimeMetricUseCase(metricRepo, metricSender)

			err := usecase.SendRunTimeMetric()
			if err != nil {
				assert.EqualError(t, d.err, err.Error())
			}
			assert.Equal(t, true, metricRepo.HasCalled)
			assert.Equal(t, d.err, err)
		})
	}
}
