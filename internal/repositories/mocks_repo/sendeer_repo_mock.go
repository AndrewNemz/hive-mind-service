package mocksrepo

import "hiv_mind/internal/entities"

type MockSenderRepo struct {
	HasCalled bool
	Metrics   []entities.Metrics
	Err       error
}

func (msr *MockSenderRepo) SendMetrics(metrics []entities.Metrics) error {
	msr.HasCalled = true
	msr.Metrics = metrics
	return msr.Err
}
