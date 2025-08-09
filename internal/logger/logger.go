package logger

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds logger configuration
type Config struct {
	LogFile    string `envconfig:"LOG_FILE" default:"logs/gimbal.log"`
	LogLevel   string `envconfig:"LOG_LEVEL" default:"DEBUG"`
	ConsoleOut bool   `envconfig:"LOG_CONSOLE_OUT" default:"true"`
	FileOut    bool   `envconfig:"LOG_FILE_OUT" default:"true"`
}

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
	file     *os.File
}

// NewWithConfig creates a new logger instance with custom configuration
func NewWithConfig(config *Config) (*Logger, error) {
	// If no config provided, create default config
	if config == nil {
		config = &Config{}
		// Use envconfig to load defaults
		if err := envconfig.Process("", config); err != nil {
			// Fallback to hardcoded defaults if envconfig fails
			config = &Config{
				LogFile:    "logs/gimbal.log",
				LogLevel:   "DEBUG",
				ConsoleOut: true,
				FileOut:    true,
			}
		}
	}

	level, err := zapcore.ParseLevel(config.LogLevel)
	if err != nil {
		level = zapcore.DebugLevel
	}

	cores, logFile, err := createLoggerCores(config, level)
	if err != nil {
		return nil, err
	}

	zapLogger := createZapLogger(cores)
	logger := createLogger(zapLogger, logFile)

	// Log initial message
	logger.logInitialization(config)

	return logger, nil
}

// createLoggerCores creates and configures the logger cores
func createLoggerCores(config *Config, level zapcore.Level) ([]zapcore.Core, *os.File, error) {
	var cores []zapcore.Core
	var logFile *os.File
	var err error

	// Add console core if enabled
	if config.ConsoleOut {
		consoleCore := createConsoleCore(level)
		cores = append(cores, consoleCore)
	}

	// Add file core if enabled
	if config.FileOut && config.LogFile != "" {
		logFile, err = createLogFile(config.LogFile)
		if err != nil {
			return nil, nil, err
		}

		fileCore := createFileCore(logFile, level)
		cores = append(cores, fileCore)
	}

	// If no cores configured, default to console
	if len(cores) == 0 {
		consoleCore := createConsoleCore(level)
		cores = append(cores, consoleCore)
	}

	return cores, logFile, nil
}

// createConsoleCore creates a console output core
func createConsoleCore(level zapcore.Level) zapcore.Core {
	consoleEncoderConfig := zapcore.EncoderConfig{
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

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		zapcore.AddSync(&syncWriter{Writer: os.Stdout}),
		level,
	)
}

// createFileCore creates a file output core
func createFileCore(logFile *os.File, level zapcore.Level) zapcore.Core {
	fileEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(fileEncoderConfig),
		zapcore.AddSync(&syncWriter{Writer: logFile}),
		level,
	)
}

// createLogFile creates and opens the log file
func createLogFile(logFilePath string) (*os.File, error) {
	// Ensure log directory exists
	logDir := filepath.Dir(logFilePath)
	if mkdirErr := os.MkdirAll(logDir, 0o755); mkdirErr != nil {
		return nil, mkdirErr
	}

	// Open log file (create if doesn't exist, append if exists)
	return os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
}

// createZapLogger creates the zap logger with cores
func createZapLogger(cores []zapcore.Core) *zap.Logger {
	// Create multi-core
	core := zapcore.NewTee(cores...)

	// Create logger with development options
	return zap.New(core,
		zap.Development(),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}

// createLogger creates the Logger wrapper
func createLogger(zapLogger *zap.Logger, logFile *os.File) *Logger {
	return &Logger{
		Logger:   zapLogger,
		lastLogs: make(map[string]any),
		file:     logFile,
	}
}

// logInitialization logs the initial logger setup message
func (l *Logger) logInitialization(config *Config) {
	l.Info("Logger initialized",
		"log_file", config.LogFile,
		"log_level", config.LogLevel,
		"console_output", config.ConsoleOut,
		"file_output", config.FileOut)
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

// DebugContext logs a debug message with context
func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...any) {
	if !l.shouldLog(msg, fields...) {
		return
	}
	l.Logger.Debug(msg, toZapFields(fields...)...)
}

// InfoContext logs an info message with context
func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...any) {
	if !l.shouldLog(msg, fields...) {
		return
	}
	l.Logger.Info(msg, toZapFields(fields...)...)
}

// WarnContext logs a warning message with context
func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...any) {
	if !l.shouldLog(msg, fields...) {
		return
	}
	l.Logger.Warn(msg, toZapFields(fields...)...)
}

// ErrorContext logs an error message with context
func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...any) {
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

	return compareByKind(va, vb)
}

// compareByKind delegates comparison based on reflect.Kind
func compareByKind(va, vb reflect.Value) bool {
	switch va.Kind() {
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128, reflect.String, reflect.UnsafePointer:
		return comparePrimitive(va, vb)
	case reflect.Slice, reflect.Array:
		return equalSlicesOrArrays(va, vb)
	case reflect.Map:
		return equalMaps(va, vb)
	case reflect.Struct:
		return equalStructs(va, vb)
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Pointer:
		return comparePrimitive(va, vb)
	default:
		return false
	}
}

// comparePrimitive compares primitive types using direct interface comparison
func comparePrimitive(va, vb reflect.Value) bool {
	return va.Interface() == vb.Interface()
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
	// Sync the underlying zap logger
	if err := l.Logger.Sync(); err != nil {
		return err
	}

	// Close the log file if it exists
	if l.file != nil {
		return l.file.Close()
	}

	return nil
}
