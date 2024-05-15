package logging

import (
	"log/slog"
	"os"
)

func NewLogger(service string, env string) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	child := logger.With(slog.Group("service_info", slog.String("env", env), slog.String("service", service)))
	return child
}

func LoggerFor(parent *slog.Logger, child string) *slog.Logger {
	return parent.With("context_source", child)
}
