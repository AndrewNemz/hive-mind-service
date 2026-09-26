package app

import (
	"hiv_mind/internal/repositories"
	psqlmetricsrepo "hiv_mind/internal/repositories/psql_metrics_repo"
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
	storeInterval int, storageFile string, db *postgresql.PostgreSQL, postgresDSN string,
) *ServiceProvider {
	storage := NewStorage(postgresDSN, db)
	metricUseCase := usecases.NewMetricUseCase(storage)
	return &ServiceProvider{
		Storage:       storage,
		MetricUseCase: metricUseCase,
		StoreInterval: storeInterval,
		StorageFile:   storageFile,
		DB:            db,
	}
}

func NewStorage(postgresDSN string, db *postgresql.PostgreSQL) repositories.IMetricStoragerRepo {
	if postgresDSN != "" {
		storage := psqlmetricsrepo.NewPostgresStorage(db)
		return storage
	}
	return repositories.NewMemStorage()
}
