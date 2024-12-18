package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger struct {
	Log *zap.Logger
}

func NewLogger() *Logger {
	l := &Logger{Log: zap.NewNop()}
	if err := l.initialize("info"); err != nil {
		log.Fatalf("Error initialize zap logger: %v", err)
	}
	return l
}

func (l *Logger) initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	l.Log = zl
	return nil
}
