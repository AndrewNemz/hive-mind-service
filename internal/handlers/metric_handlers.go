package handlers

import (
	"errors"
	"fmt"
	"hiv_mind/internal/app"
	"hiv_mind/internal/entities"
	repoerrors "hiv_mind/internal/errors"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MetricHandler struct {
	serviceProvider *app.ServiceProvider
	Template        *template.Template
}

func NewMetricHandler(sp *app.ServiceProvider, templatesDir string) (*MetricHandler, error) {

	tmplPath := filepath.Join(templatesDir, "metrics.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, err
	}
	return &MetricHandler{
		serviceProvider: sp,
		Template:        tmpl,
	}, nil
}

func (mh *MetricHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metric, err := ValidateParams(w, r)
	if err != nil {
		return
	}

	err = mh.serviceProvider.MetricUseCase.CollectAndStoreMetric(metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(mh.serviceProvider.Storage)
	w.WriteHeader(http.StatusOK)
	return
}

func ValidateParams(w http.ResponseWriter, r *http.Request) (entities.Metrics, error) {

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if metricType == "" || metricName == "" || metricValue == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("Не переданы обязательные параметры!")
		return entities.Metrics{}, fmt.Errorf("Неккоректные параметры запроса!")
	}

	metric := &entities.Metrics{Name: metricName}

	switch metricType {
	case entities.GaugeType:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return entities.Metrics{}, fmt.Errorf("Ошибка обработки метрики!")
		}
		metric.Value = value
		metric.Type = metricType
	case entities.CounterType:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return entities.Metrics{}, fmt.Errorf("Ошибка обработки метрики!")
		}
		metric.Value = float64(value)
		metric.Type = metricType
	default:
		w.WriteHeader(http.StatusBadRequest)
		return entities.Metrics{}, fmt.Errorf("Некорректный тип метрики")
	}

	return *metric, nil
}

func (mh *MetricHandler) Value(w http.ResponseWriter, r *http.Request) {

	metricName := chi.URLParam(r, "name")
	metricType := chi.URLParam(r, "type")

	if metricType != entities.CounterType && metricType != entities.GaugeType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := mh.serviceProvider.MetricUseCase.GetMetricByTypeAndName(metricType, metricName)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFoundMetric) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("%v", metric.Value)))

	return
}

func (mh *MetricHandler) Root(w http.ResponseWriter, r *http.Request) {

	metrics := mh.serviceProvider.MetricUseCase.GetAllMetrics()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := mh.Template.Execute(w, metrics); err != nil {
		http.Error(w, "Ошибка генерации HTML", http.StatusInternalServerError)
		return
	}
}
