package usecases

import (
	"fmt"
	"hiv_mind/internal/entities"
	"hiv_mind/internal/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricUseCase(t *testing.T) {
	data := []struct {
		tName string
		metic entities.Metrics
		err   error
	}{
		{
			tName: "Test_Method_CollectAndStoreMetric_Counter_When_OK",
			metic: entities.Metrics{Type: entities.GaugeType, Name: "Alloc", Value: 12.0},
			err:   nil,
		},
		{
			tName: "Test_Method_CollectAndStoreMetric_Gauge_When_OK",
			metic: entities.Metrics{Type: entities.CounterType, Name: "Alloc", Value: 10},
			err:   nil,
		},
		{
			tName: "Test_Method_CollectAndStoreMetric_When_Invalid_Type",
			metic: entities.Metrics{Type: "InvalidType", Name: "Alloc", Value: 10},
			err:   fmt.Errorf("В репозиторий передан неожиданный формат метрики!"),
		},
	}

	for _, d := range data {
		t.Run(d.tName, func(t *testing.T) {
			repo := repositories.NewMemStorage()
			usecase := NewMetricUseCase(repo)

			err := usecase.CollectAndStoreMetric(d.metic)

			if err != nil {
				assert.EqualError(t, d.err, err.Error())
			}
			assert.Equal(t, d.err, err)
		})
	}
}
