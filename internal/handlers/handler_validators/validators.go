package handlervalidators

import (
	"encoding/json"
	"fmt"
	"hiv_mind/internal/entities"
	"net/http"

	"go.uber.org/zap"
)

func ValidateParams(w http.ResponseWriter, r *http.Request, lg *zap.Logger) (*entities.Metrics, error) {

	var metric entities.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		lg.Info("Failed to decode JSON body", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	if metric.MType == "" || metric.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		lg.Info("Не переданы обязательные параметры!")
		return nil, fmt.Errorf("Неккоректные параметры запроса!")
	}
	if metric.Value == nil && metric.Delta == nil {
		w.WriteHeader(http.StatusNotFound)
		lg.Info("Не переданы обязательные параметры!")
		return nil, fmt.Errorf("Неккоректные параметры запроса!")
	}
	if metric.MType != entities.GaugeType && metric.MType != entities.CounterType {
		w.WriteHeader(http.StatusBadRequest)
		lg.Info("Не переданы обязательные параметры!")
		return nil, fmt.Errorf("Неккоректные параметры запроса!")
	}

	return &metric, nil
}
