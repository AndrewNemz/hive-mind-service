package repositories

import (
	"fmt"
	"hiv_mind/internal/entities"

	"github.com/go-resty/resty/v2"
)

type MetricSender struct {
	client *resty.Client
}

func NewMetricSender() *MetricSender {
	return &MetricSender{
		client: resty.New().SetBaseURL("http://localhost:8080/update/"),
	}
}

func (ms *MetricSender) SendMetrics(metrics []entities.Metrics) error {

	for _, m := range metrics {
		url := fmt.Sprintf("%s/%s/%g", m.Type, m.Name, m.Value)
		fmt.Println(url)

		response, err := ms.client.R().SetHeader(`Content-Type`, "text/plain").Post(url)
		if err != nil {
			fmt.Printf("Ошибка отправки метрики %s: %v\n", m.Name, err)
			continue
		}

		fmt.Println(response)
	}

	return nil
}
