package repositories

import (
	"encoding/json"
	"fmt"
	"hiv_mind/internal/entities"
	compressdata "hiv_mind/pkg/compress_data"
	"hiv_mind/pkg/logger"

	"github.com/go-resty/resty/v2"
)

type MetricSender struct {
	client  *resty.Client
	Adresss string
}

func NewMetricSender(adresss string) *MetricSender {
	return &MetricSender{
		client: resty.New().
			SetBaseURL(fmt.Sprintf("http://%s/update", adresss)).
			SetHeader("Content-Encoding", "gzip").
			SetHeader(`Content-Type`, "application/json"),
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

		compressData, err := compressdata.Compress(jsonData)
		if err != nil {
			lg.Sugar().Errorf("Ошибка сжатия данных %s: %v\n", m.ID, err)
			continue
		}

		_, err = ms.client.R().SetBody(compressData).Post("/")
		if err != nil {
			lg.Sugar().Errorf("Ошибка отправки метрики %s: %v\n", m.ID, err)
			continue
		}
	}

	return nil
}
