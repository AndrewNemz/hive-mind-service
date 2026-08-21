package handlers

import (
	"errors"
	"fmt"
	"hiv_mind/internal/app"
	"hiv_mind/internal/entities"
	repoerrors "hiv_mind/internal/errors"
	"hiv_mind/pkg/logger"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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
	lg := logger.Get()

	if r.Method != http.MethodPost {
		lg.Info(fmt.Sprintf("Request Method %s is not Allowded", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metric, err := ValidateParams(w, r, lg)
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

func ValidateParams(w http.ResponseWriter, r *http.Request, lg *zap.Logger) (entities.Metrics, error) {

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if metricType == "" || metricName == "" || metricValue == "" {
		w.WriteHeader(http.StatusNotFound)
		lg.Info("Не переданы обязательные параметры!")
		return entities.Metrics{}, fmt.Errorf("Неккоректные параметры запроса!")
	}

	metric := &entities.Metrics{Name: metricName}

	switch metricType {
	case entities.GaugeType:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			lg.Info("Ошибка обработки метрики!")
			w.WriteHeader(http.StatusBadRequest)
			return entities.Metrics{}, fmt.Errorf("Ошибка обработки метрики!")
		}
		metric.Value = value
		metric.Type = metricType
	case entities.CounterType:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			lg.Info("Ошибка обработки метрики!")
			w.WriteHeader(http.StatusBadRequest)
			return entities.Metrics{}, fmt.Errorf("Ошибка обработки метрики!")
		}
		metric.Value = float64(value)
		metric.Type = metricType
	default:
		lg.Info("Некорректный тип метрики")
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

	lg := logger.Get()
	metrics := mh.serviceProvider.MetricUseCase.GetAllMetrics()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := mh.Template.Execute(w, metrics); err != nil {
		lg.Info("Ошибка генерации HTML")
		http.Error(w, "Ошибка генерации HTML", http.StatusInternalServerError)
		return
	}
}
