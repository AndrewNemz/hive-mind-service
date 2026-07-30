package handlers

import (
	"fmt"
	"hiv_mind/internal/app"
	"hiv_mind/internal/entities"
	"net/http"
	"strconv"
	"strings"
)

type MetricHandler struct {
	serviceProvider *app.ServiceProvider
}

func NewMetricHandler(sp *app.ServiceProvider) *MetricHandler {
	return &MetricHandler{
		serviceProvider: sp,
	}
}

func (mh *MetricHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	metricType := params[2]
	metricName := params[3]

	var metricValue float64
	var metricValueF float64
	var metricValueI int64

	if metricType == "gauge" {
		metricValueF, err = strconv.ParseFloat(params[4], 64)
		metricValue = metricValueF
	} else {
		metricValueI, err = strconv.ParseInt(params[4], 10, 64)
		metricValue = float64(metricValueI)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metric := entities.Metrics{
		Type:  metricType,
		Name:  metricName,
		Value: metricValue,
	}
	fmt.Println(metric)
	err = mh.serviceProvider.MetricUseCase.CollectAndStoreMetric(metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(mh.serviceProvider.Storage)
	w.WriteHeader(http.StatusOK)
	return
}

func ValidateParams(params []string, w http.ResponseWriter) error {
	if len(params) < 5 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("Не переданы обязательные параметры!")
		return fmt.Errorf("Неккоректные параметры запроса!")
	}
	typeM, _, _ := params[2], params[3], params[4]

	switch typeM {
	case "gauge":
	case "counter":
	default:
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("Некорректный тип метрики")
	}

	return nil
}
