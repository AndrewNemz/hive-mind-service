package repoerrors

import "errors"

var (
	ErrNotFoundMetric = errors.New("Метрика не найдена в хранилище!")
)
