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
	err := ValidateParams(params, w)
	if err != nil {
		return
	}

	metricType := params[0]
	metricName := params[1]

	var valInt int64
	var valFloat float64
	if metricType == "gauge" {
		valFloat, err = strconv.ParseFloat(params[2], 64)
	} else {
		valInt, err = strconv.ParseInt(params[2], 10, 64)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, ok := ms.Storage[metricName]; !ok {
		ms.Storage[metricName] = &entities.Metrics{}
	}

	switch strings.ToLower(metricType) {
	case "gauge":
		ms.Storage[metricName].Gauge = valFloat
		fmt.Println(ms.Storage[metricName])
		w.WriteHeader(http.StatusOK)
		return
	case "counter":
		ms.Storage[metricName].Counter += valInt
		fmt.Println(ms.Storage[metricName])
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.Write([]byte(fmt.Sprintf("Тип метрики %s некорректен", metricType)))
		return
	}

}

func ValidateParams(params []string, w http.ResponseWriter) error {
	fmt.Println("Переданные параметры", params)
	if len(params) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("Не переданы обязательные параметры!")
		return fmt.Errorf("Неккоректные параметры запроса!")
	}
	typeM, _, _ := params[0], params[1], params[2]

	switch typeM {
	case "gauge":
	case "counter":
	default:
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("Некорректный тип метрики")
	}

	return nil
}
