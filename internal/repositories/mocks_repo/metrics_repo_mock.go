package mocksrepo

import "hiv_mind/internal/entities"

type MockStorageRepo struct {
	HasCalled bool
	Metric    entities.Metrics
	Metrics   []entities.Metrics
	Err       error
}

func (msr *MockStorageRepo) StoreMetric(m entities.Metrics) error {
	msr.HasCalled = true
	msr.Metric = m
	return msr.Err
}

func (msr *MockStorageRepo) StoreMetricSlice(metrics []entities.Metrics) error {
	msr.HasCalled = true
	msr.Metrics = metrics
	return msr.Err
}

func (msr *MockStorageRepo) GetAllMetrics() []entities.Metrics {
	msr.HasCalled = true
	return msr.Metrics
}

func (msr *MockStorageRepo) GetMetricByTypeAndName(mType, mName string) (entities.Metrics, error) {
	msr.HasCalled = true
	return entities.Metrics{Name: mName, Type: mType, Value: 100}, msr.Err
}
