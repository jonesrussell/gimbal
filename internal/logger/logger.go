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
	// Create stdout syncer
	stdout := zapcore.Lock(os.Stdout)

	// Get log level from environment variable, default to debug
	level := zapcore.DebugLevel
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

	// Create encoder config with more readable format
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core with console output
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(consoleEncoder, stdout, level)

	// Create logger with fields and options
	GlobalLogger = zap.New(core,
		zap.AddCaller(),
		zap.Development(),
	)

	// Log initial message to verify logger is working
	GlobalLogger.Info("Logger initialized",
		zap.String("level", level.String()),
		zap.String("go_version", runtime.Version()),
	)
}

// Sync flushes any buffered log entries
func Sync() error {
	return GlobalLogger.Sync()
}
