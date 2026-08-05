package repositories

import (
	"fmt"
	"hiv_mind/internal/entities"

	"github.com/go-resty/resty/v2"
)

type MetricSender struct {
	client  *resty.Client
	Adresss string
}

func NewMetricSender(adresss string) *MetricSender {
	return &MetricSender{
		client:  resty.New().SetBaseURL(fmt.Sprintf("http://%s/update/", adresss)),
		Adresss: adresss,
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
