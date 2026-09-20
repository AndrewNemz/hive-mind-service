package app

import (
	"hiv_mind/internal/repositories"
	usecases "hiv_mind/internal/use_cases"
	"hiv_mind/pkg/postgresql"
)

type ServiceProvider struct {
	Storage       repositories.IMetricStoragerRepo
	MetricUseCase usecases.IMetricsUseCase
	StoreInterval int
	StorageFile   string
	DB            *postgresql.PostgreSQL
}

func NewServiceProvider(
	storeInterval int, storageFile string, db *postgresql.PostgreSQL,
) *ServiceProvider {
	storage := repositories.NewMemStorage()
	metricUseCase := usecases.NewMetricUseCase(storage)
	return &ServiceProvider{
		Storage:       storage,
		MetricUseCase: metricUseCase,
		StoreInterval: storeInterval,
		StorageFile:   storageFile,
		DB:            db,
	}
}
