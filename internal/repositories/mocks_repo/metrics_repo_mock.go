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
	var delta int64 = 10
	return entities.Metrics{ID: "123", MType: entities.CounterType, Delta: &delta}, msr.Err
}

func (msr *MockStorageRepo) LoadMetricFromFile(filename string) error {
	msr.HasCalled = true
	return nil
}

func (msr *MockStorageRepo) SaveToFile(filename string) error {
	return nil
}
