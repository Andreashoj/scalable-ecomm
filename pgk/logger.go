package pgk

import (
	"log/slog"
	"os"
)

type Logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
}

func NewLogger() Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
