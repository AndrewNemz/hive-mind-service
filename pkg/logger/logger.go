package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}

func Get() *zap.Logger {
	if Log == nil {
		return zap.NewNop()
	}
	return Log
}

func SugaredLogger() *zap.SugaredLogger {
	return Get().Sugar()
}

func DesugaredLogger(sugar *zap.SugaredLogger) *zap.Logger {
	return sugar.Desugar()
}
