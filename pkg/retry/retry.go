package retry

import (
	"context"
	"errors"
	"hiv_mind/pkg/logger"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type CustomBackOff struct {
	Ctx         context.Context
	MaxAttempts int
	Durations   []time.Duration
}

func NewCustomBackOff() *CustomBackOff {
	return &CustomBackOff{
		Ctx:         context.TODO(),
		MaxAttempts: 3,                                                                  // Вынести в .env
		Durations:   []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}, // вынести в .env
	}
}

func IsRetriableDBError(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return len(pgErr.Code) >= 2 && pgErr.Code[:2] == "08"
	}

	return false
}

func IsRetriableHTTPError(statusCode int, err error) bool {
	if err != nil {
		var netErr net.Error
		return errors.As(err, &netErr)
	}

	return statusCode >= 500
}

func (cbo *CustomBackOff) Retry(f func() error, isRetriable func(error) bool) error {

	var err error
	lg := logger.Get()

	for i := 0; i < cbo.MaxAttempts; i++ {
		err = f()

		if err == nil {
			break
		}
		if !isRetriable(err) {
			lg.Error("Невосстановимая ошибка, повторные попытки отменены", zap.Error(err))
			return err // Прерываем цикл немедленно
		}
		if i == cbo.MaxAttempts-1 {
			break
		}

		lg.Info("Функция отработала с ошибкой! Вызван механизм Rerty - повтор через %d", zap.Duration("время, с", cbo.Durations[i]))
		time.Sleep(cbo.Durations[i])

		// select {
		// case <-ctx.Done():
		// 	return ctx.Err()
		// case <-time.After(delay):
		// 	// Задержка прошла, идем на следующий круг цикла
		// }
	}

	return err
}
