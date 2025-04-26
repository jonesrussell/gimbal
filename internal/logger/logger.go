package logger

import (
	"io"
	"os"
	"reflect"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// syncWriter wraps an io.Writer to ignore sync errors for stdout/stderr
type syncWriter struct {
	io.Writer
}

func (w *syncWriter) Sync() error {
	// Ignore sync errors for stdout/stderr
	return nil
}

// Logger wraps a zap.Logger to implement common.Logger
type Logger struct {
	*zap.Logger
	lastLogs map[string]interface{}
	mu       sync.RWMutex
}

// New creates a new logger instance
func New() (*Logger, error) {
	// Create encoder for console output
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core with console output using our custom writer
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(&syncWriter{os.Stdout}),
		zapcore.DebugLevel,
	)

	// Build the logger
	zapLogger := zap.New(core,
		zap.Development(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	logger := &Logger{
		Logger:   zapLogger,
		lastLogs: make(map[string]interface{}),
	}

	// Log initial message
	logger.Info("Logger initialized")

	return logger, nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...any) {
	if !l.shouldLog(msg, fields...) {
		return
	}
	l.Logger.Debug(msg, toZapFields(fields...)...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...any) {
	if !l.shouldLog(msg, fields...) {
		return
	}
	l.Logger.Info(msg, toZapFields(fields...)...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...any) {
	if !l.shouldLog(msg, fields...) {
		return
	}
	l.Logger.Warn(msg, toZapFields(fields...)...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...any) {
	// Always log errors, don't deduplicate them
	l.Logger.Error(msg, toZapFields(fields...)...)
}

// shouldLog determines if a message should be logged based on deduplication rules
func (l *Logger) shouldLog(msg string, fields ...any) bool {
	// Don't deduplicate if there are no fields
	if len(fields) == 0 {
		return true
	}

	// Create a key from the message and first field value
	key := msg
	if len(fields) >= 2 {
		if str, ok := fields[0].(string); ok {
			key += ":" + str
		}
	}

	// Get the current value
	currentValue := fields

	l.mu.Lock()
	defer l.mu.Unlock()

	// Check if this is a duplicate
	if lastValue, exists := l.lastLogs[key]; exists {
		if equalValues(lastValue, currentValue) {
			return false
		}
	}

	// Update the last logged value
	l.lastLogs[key] = currentValue
	return true
}

// equalValues compares two sets of field values for equality
func equalValues(a, b interface{}) bool {
	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Convert to reflect.Value for type-safe comparison
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	// If types don't match, they're not equal
	if va.Type() != vb.Type() {
		return false
	}

	// Handle different types appropriately
	switch va.Kind() {
	case reflect.Slice:
		if va.Len() != vb.Len() {
			return false
		}
		for i := 0; i < va.Len(); i++ {
			if !equalValues(va.Index(i).Interface(), vb.Index(i).Interface()) {
				return false
			}
		}
		return true
	case reflect.Map:
		if va.Len() != vb.Len() {
			return false
		}
		for _, k := range va.MapKeys() {
			if !equalValues(va.MapIndex(k).Interface(), vb.MapIndex(k).Interface()) {
				return false
			}
		}
		return true
	case reflect.Struct:
		for i := 0; i < va.NumField(); i++ {
			if !equalValues(va.Field(i).Interface(), vb.Field(i).Interface()) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(a, b)
	}
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

// Sync ensures all buffered logs are written
func (l *Logger) Sync() error {
	// Ignore sync errors for stdout
	return nil
}
