package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"hiv_mind/internal/app"
	"hiv_mind/internal/handlers"
	"hiv_mind/pkg/logger"
	"hiv_mind/pkg/middleware"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var adress string
var storeInterval int
var storageFile string
var restoreValues bool

func init() {
	flag.StringVar(&adress, "a", "localhost:8080", "отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080)")
	flag.IntVar(&storeInterval, "i", 300, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&storageFile, "f", "/tmp/metrics-db.json", "полное имя файла, куда сохраняются текущие значения")
	flag.BoolVar(&restoreValues, "r", true, "булево значение (true/false), определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize("info"); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	lg := logger.Get()
	defer lg.Sync()

	flag.Parse()
	err := GetEnvironment()
	if err != nil {
		return err
	}
	lg.Info("Старт сервиса",
		zap.String("address", adress),
		zap.Int("store_interval_sec", storeInterval),
		zap.String("storage_file", storageFile),
		zap.Bool("restore", restoreValues),
	)

	r := chi.NewRouter()
	serviceProvider := app.NewServiceProvider(storeInterval, storageFile)
	metricHandler, err := handlers.NewMetricHandler(serviceProvider, "./templates")
	if err != nil {
		lg.Error("не удалось загрузить шаблоны", zap.Error(err))
		return err
	}

	if restoreValues && storageFile != "" {
		if err := serviceProvider.Storage.LoadMetricFromFile(storageFile); err != nil {
			lg.Error("не удалось загрузить метрики из файла", zap.Error(err))
		} else {
			lg.Info("Метрики успешно восстановлены из файла")
		}
	}

	// middlewares
	r.Use(middleware.GzipMiddleware)
	r.Use(middleware.Logging)

	// Routes
	r.Get("/", metricHandler.Root)
	r.Post("/value/", metricHandler.Value)
	r.Post("/update/", metricHandler.Update)

	server := &http.Server{
		Addr:    adress,
		Handler: r,
	}
	// Контекст, который отменится при получении SIGINT или SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if storeInterval > 0 && storageFile != "" {
		go func() {
			ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err := serviceProvider.Storage.SaveToFile(storageFile); err != nil {
						lg.Error("Ошибка периодического сохранения", zap.Error(err))
					}
				}
			}
		}()
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		lg.Sugar().Infof("сервер запущен на %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Sugar().Fatalf("ошибка сервера: %v", err)
		}
	}()

	// Ждём сигнала
	<-ctx.Done()
	lg.Info("получен сигнал завершения, начинаем shutdown")

	if storageFile != "" {
		lg.Info("Сохранение накопленных данных на диск...")
		if err := serviceProvider.Storage.SaveToFile(storageFile); err != nil {
			lg.Error("Ошибка при финальном сохранении", zap.Error(err))
		} else {
			lg.Info("Данные успешно сохранены")
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		lg.Sugar().Infof("ошибка при shutdown: %v", err)
	}

	lg.Info("сервер остановлен")

	return nil
}

func GetEnvironment() error {
	lg := logger.Get()

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		val, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			lg.Error("Ошибка конвертации флага i", zap.Error(err))
			return err
		}
		storeInterval = val
	}

	if envStorageFile := os.Getenv("FILE_STORAGE_PATH"); envStorageFile != "" {
		storageFile = envStorageFile
	}

	if envRestoreValues := os.Getenv("RESTORE"); envRestoreValues != "" {
		val, err := strconv.ParseBool(envRestoreValues)
		if err != nil {
			lg.Error("Ошибка конвертации флага r", zap.Error(err))
			return err
		}
		restoreValues = val
	}

	if envAdress := os.Getenv("ADDRESS"); envAdress != "" {
		adress = envAdress
	}

	return nil
}
