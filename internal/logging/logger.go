package logging

import (
	"context"
	"log/slog"
	"os"
)

type LoggerConfig struct {
	Level slog.Level
}

func New(component string, cfg LoggerConfig) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Key = "timestamp"
			}
			return a
		},
	})
	return slog.New(handler).With("component", component)
}
func СonfigureLogger(debug, info, error bool) *slog.Logger {
	var level slog.Level

	switch {
	case debug:
		level = slog.LevelDebug
	case info:
		level = slog.LevelInfo
	case error:
		level = slog.LevelError
	default:
		// Если флаги не указаны, то отображаться будет только критические ошибки Error
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}

// Вспомогательные методы для логирования (думал будет компактней, но как-то хз)
func LogInfo(logger *slog.Logger, msg string, attrs ...slog.Attr) {
	if logger != nil && logger.Enabled(context.Background(), slog.LevelInfo) {
		logger.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
	}
}
func LogError(logger *slog.Logger, err error, msg string, attrs ...slog.Attr) {
	if logger != nil && logger.Enabled(context.Background(), slog.LevelError) {
		attrs = append(attrs, slog.String("error", err.Error()))
		logger.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
	}
}
