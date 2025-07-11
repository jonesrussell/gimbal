package logger

import (
	"io"
	"os"
	"reflect"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// syncWriter wraps an io.Writer to make it safe for concurrent use
type syncWriter struct {
	io.Writer
	mu sync.Mutex
}

func (w *syncWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Writer.Write(p)
}

// Logger wraps zap.Logger with additional functionality
type Logger struct {
	*zap.Logger
	lastLogs map[string]any
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
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core with console output
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(&syncWriter{Writer: os.Stdout}),
		zapcore.DebugLevel,
	)

	// Create logger with development options
	zapLogger := zap.New(core,
		zap.Development(),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	logger := &Logger{
		Logger:   zapLogger,
		lastLogs: make(map[string]any),
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

// equalValues compares two values for equality
func equalValues(a, b any) bool {
	if a == nil || b == nil {
		return a == b
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Type() != vb.Type() {
		return false
	}

	switch va.Kind() {
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128, reflect.String, reflect.UnsafePointer:
		return va.Interface() == vb.Interface()
	case reflect.Slice, reflect.Array:
		return equalSlicesOrArrays(va, vb)
	case reflect.Map:
		return equalMaps(va, vb)
	case reflect.Struct:
		return equalStructs(va, vb)
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Pointer:
		return va.Interface() == vb.Interface()
	default:
		return false
	}
}

func equalSlicesOrArrays(va, vb reflect.Value) bool {
	if va.Len() != vb.Len() {
		return false
	}
	for i := 0; i < va.Len(); i++ {
		if !equalValues(va.Index(i).Interface(), vb.Index(i).Interface()) {
			return false
		}
	}
	return true
}

func equalMaps(va, vb reflect.Value) bool {
	if va.Len() != vb.Len() {
		return false
	}
	for _, k := range va.MapKeys() {
		if !equalValues(va.MapIndex(k).Interface(), vb.MapIndex(k).Interface()) {
			return false
		}
	}
	return true
}

func equalStructs(va, vb reflect.Value) bool {
	for i := 0; i < va.NumField(); i++ {
		if !equalValues(va.Field(i).Interface(), vb.Field(i).Interface()) {
			return false
		}
	}
	return true
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
