package logger

import (
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

var GlobalLogger slog.Logger

func init() {
	// Get log level from environment variable, default to info
	level := slog.LevelInfo
	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		switch strings.ToUpper(levelStr) {
		case "DEBUG":
			level = slog.LevelDebug
		case "INFO":
			level = slog.LevelInfo
		case "WARN":
			level = slog.LevelWarn
		case "ERROR":
			level = slog.LevelError
		}
	}

	// Initialize the global logger with the configured level
	GlobalLogger = NewSlogHandler(level)
}

func NewSlogHandler(level slog.Level) slog.Logger {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level, // Use the provided logging level
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Key = "UTCTime"
				a.Value = slog.TimeValue(time.Now().UTC())
			}
			return a
		},
	}).WithAttrs([]slog.Attr{
		slog.Group("app_details",
			slog.Int("pid", os.Getpid()),
			slog.String("go_version", runtime.Version()),
		),
	})

	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	return *logger
}
