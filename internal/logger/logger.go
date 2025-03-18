package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps a zap.Logger to implement common.Logger
type Logger struct {
	*zap.Logger
}

// New creates a new logger instance
func New() (*Logger, error) {
	// Create a basic console encoder
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		Logger: zapLogger,
	}

	// Log initial message
	logger.Info("Logger initialized")

	return logger, nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...any) {
	l.Logger.Debug(msg, toZapFields(fields...)...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...any) {
	l.Logger.Info(msg, toZapFields(fields...)...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...any) {
	l.Logger.Warn(msg, toZapFields(fields...)...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...any) {
	l.Logger.Error(msg, toZapFields(fields...)...)
}

// toZapFields converts interface slice to zap.Field slice
func toZapFields(fields ...any) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if !ok {
				continue
			}
			zapFields = append(zapFields, zap.Any(key, fields[i+1]))
		}
	}
	return zapFields
}
