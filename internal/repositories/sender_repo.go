package repositories

import (
	"fmt"
	"hiv_mind/internal/entities"
	"net/http"
)

type MetricSender struct{}

func NewMetricSender() *MetricSender {
	return &MetricSender{}
}

func (ms *MetricSender) SendMetrics(metrics []entities.Metrics) error {
	baseURL := "http://localhost:8080/update/"

	var url string
	client := &http.Client{}

	for _, m := range metrics {
		url = baseURL + fmt.Sprintf("%s/%s/%g", m.Type, m.Name, m.Value)
		fmt.Println(url)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			fmt.Printf("Ошибка создания запроса для %s: %v\n", m.Name, err)
			continue
		}
		request.Header.Set(`Content-Type`, "text/plain")

		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("Ошибка отправки метрики %s: %v\n", m.Name, err)
			continue
		}

		fmt.Println(response)
		response.Body.Close()
	}

	return nil
}
