package logger

import (
	"os"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var GlobalLogger *zap.Logger

func init() {
	// Get log level from environment variable, default to info
	level := zapcore.InfoLevel
	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		switch strings.ToUpper(levelStr) {
		case "DEBUG":
			level = zapcore.DebugLevel
		case "INFO":
			level = zapcore.InfoLevel
		case "WARN":
			level = zapcore.WarnLevel
		case "ERROR":
			level = zapcore.ErrorLevel
		}
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Create logger with fields
	GlobalLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).With(
		zap.Int("pid", os.Getpid()),
		zap.String("go_version", runtime.Version()),
	)
}
