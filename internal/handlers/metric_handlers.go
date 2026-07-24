package handlers

import (
	"fmt"
	"hiv_mind/internal/entities"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	Storage map[string]*entities.Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Storage: map[string]*entities.Metrics{},
	}
}

func (ms *MemStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	params := strings.Split(path, "/")
	if len(params) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metricType := params[0]
	metricName := params[1]
	metricVal, err := strconv.Atoi(params[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed paramas: type - %s, name - %s, val %d \n", metricType, metricName, metricVal)

	if _, ok := ms.Storage[metricName]; !ok {
		ms.Storage[metricName] = &entities.Metrics{}
	}

	switch metricType {
	case "gauge":
		ms.Storage[metricName].Gauge = float64(metricVal)
		fmt.Println(ms.Storage[metricName])
		w.WriteHeader(http.StatusOK)
		return
	case "counter":
		ms.Storage[metricName].Counter += int64(metricVal)
		fmt.Println(ms.Storage[metricName])
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.Write([]byte(fmt.Sprintf("Тип метрики %s некорректен", metricType)))
		return
	}

}
