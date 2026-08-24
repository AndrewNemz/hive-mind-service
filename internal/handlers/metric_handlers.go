package handlers

import (
	"encoding/json"
	"errors"
	"hiv_mind/internal/app"
	"hiv_mind/internal/entities"
	repoerrors "hiv_mind/internal/errors"
	handlervalidators "hiv_mind/internal/handlers/handler_validators"
	"hiv_mind/pkg/logger"
	"html/template"
	"net/http"
	"path/filepath"

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
		lg.Info("Request method not allowed", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metric, err := handlervalidators.ValidateParams(w, r, lg)
	if err != nil {
		return
	}

	err = mh.serviceProvider.MetricUseCase.CollectAndStoreMetric(metric)
	if err != nil {
		lg.Error("Failed to store metric", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (mh *MetricHandler) Value(w http.ResponseWriter, r *http.Request) {

	lg := logger.Get()

	var metric entities.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		lg.Info("Failed to decode JSON body", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lg.Sugar().Infof("Получена метрика с типом - %s, именем - %s", metric.MType, metric.ID)
	if metric.MType != entities.CounterType && metric.MType != entities.GaugeType {
		lg.Info("Invalid metric type")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := mh.serviceProvider.MetricUseCase.GetMetricByTypeAndName(&metric)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFoundMetric) {
			lg.Info("Metric not found")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Метрика не найдена"}`))
			return
		}
		lg.Error("Failed to get metric", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metric); err != nil {
		lg.Error("Failed to encode metric to JSON", zap.Error(err))
		return
	}
}

func (mh *MetricHandler) Root(w http.ResponseWriter, r *http.Request) {

	lg := logger.Get()
	metrics := mh.serviceProvider.MetricUseCase.GetAllMetrics()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := mh.Template.Execute(w, metrics); err != nil {
		lg.Error("Ошибка генерации HTML", zap.Error(err))
		http.Error(w, "Ошибка генерации HTML", http.StatusInternalServerError)
		return
	}
}
