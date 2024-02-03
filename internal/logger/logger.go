package logger

import (
	"log/slog"
	"os"
	"runtime"
	"time"
)

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
