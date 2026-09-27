package repositories

import (
	"encoding/json"
	"fmt"
	"hiv_mind/internal/entities"
	compressdata "hiv_mind/pkg/compress_data"
	"hiv_mind/pkg/logger"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type MetricSender struct {
	client  *resty.Client
	Adresss string
}

func NewMetricSender(adresss string) *MetricSender {
	return &MetricSender{
		client: resty.New().
			SetBaseURL(fmt.Sprintf("http://%s/updates/", adresss)).
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

func (ms *MetricSender) SendBatchMetrics(metrics []entities.Metrics) error {
	lg := logger.Get()

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		lg.Error("Ошибка декордирования", zap.Error(err))
		return err
	}

	compressData, err := compressdata.Compress(jsonData)
	if err != nil {
		lg.Error("Ошибка сжатия данных", zap.Error(err))
		return err
	}

	_, err = ms.client.R().SetBody(compressData).Post("/")
	if err != nil {
		lg.Error("", zap.Error(err))
		return err
	}

	return nil
}
