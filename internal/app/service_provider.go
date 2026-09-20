package app

import (
	"hiv_mind/internal/repositories"
	usecases "hiv_mind/internal/use_cases"
)

type ServiceProvider struct {
	Storage       repositories.IMetricStoragerRepo
	MetricUseCase usecases.IMetricsUseCase
	StoreInterval int
	StorageFile   string
}

func NewServiceProvider(
	storeInterval int, storageFile string,
) *ServiceProvider {
	storage := repositories.NewMemStorage()
	metricUseCase := usecases.NewMetricUseCase(storage)
	return &ServiceProvider{
		Storage:       storage,
		MetricUseCase: metricUseCase,
		StoreInterval: storeInterval,
		StorageFile:   storageFile,
	}
}
