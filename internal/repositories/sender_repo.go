package repositories

import (
	"encoding/json"
	"fmt"
	"hiv_mind/internal/entities"
	"hiv_mind/pkg/logger"

	"github.com/go-resty/resty/v2"
)

type MetricSender struct {
	client  *resty.Client
	Adresss string
}

func NewMetricSender(adresss string) *MetricSender {
	return &MetricSender{
		client:  resty.New().SetBaseURL(fmt.Sprintf("http://%s/update", adresss)).SetHeader(`Content-Type`, "application/json"),
		Adresss: adresss,
	}
}

func (ms *MetricSender) SendMetrics(metrics []entities.Metrics) error {

	lg := logger.Get()
	for _, m := range metrics {

		jsonData, err := json.Marshal(m)
		if err != nil {
			lg.Sugar().Errorf("Ошибка декордирования %s: %v\n", m.ID, err)
			continue
		}

		_, err = ms.client.R().SetBody(jsonData).Post("/")
		if err != nil {
			lg.Sugar().Errorf("Ошибка отправки метрики %s: %v\n", m.ID, err)
			continue
		}
	}

	return nil
}
