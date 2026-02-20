package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kelseyhightower/envconfig"
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

// multiHandler forwards log records to multiple handlers (e.g. console + file).
type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle forwards the record to all enabled handlers. Implements slog.Handler.
//
//nolint:gocritic // slog.Handler.Handle requires Record by value
func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	out := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		out[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: out}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	out := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		out[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: out}
}

// Logger wraps slog with deduplication and optional file for Sync.
type Logger struct {
	slog     *slog.Logger
	lastLogs map[string]any
	mu       sync.RWMutex
	file     *os.File
}

// parseLevel maps a string (e.g. "DEBUG", "INFO") to slog.Level. Defaults to Debug.
func parseLevel(s string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

// NewWithConfig creates a new logger instance with custom configuration.
func NewWithConfig(config *Config) (*Logger, error) {
	if config == nil {
		config = &Config{}
		if err := envconfig.Process("", config); err != nil {
			config = &Config{
				LogFile:    "logs/gimbal.log",
				LogLevel:   "DEBUG",
				ConsoleOut: true,
				FileOut:    true,
			}
		}
	}

	level := parseLevel(config.LogLevel)
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	var handlers []slog.Handler

	if config.ConsoleOut {
		handlers = append(handlers, slog.NewTextHandler(&syncWriter{Writer: os.Stdout}, opts))
	}

	var logFile *os.File
	if config.FileOut && config.LogFile != "" {
		var err error
		logFile, err = createLogFile(config.LogFile)
		if err != nil {
			return nil, err
		}
		handlers = append(handlers, slog.NewJSONHandler(&syncWriter{Writer: logFile}, opts))
	}

	if len(handlers) == 0 {
		handlers = append(handlers, slog.NewTextHandler(&syncWriter{Writer: os.Stdout}, opts))
	}

	var rootHandler slog.Handler = &multiHandler{handlers: handlers}
	if len(handlers) == 1 {
		rootHandler = handlers[0]
	}

	l := &Logger{
		slog:     slog.New(rootHandler),
		lastLogs: make(map[string]any),
		file:     logFile,
	}
	l.logInitialization(config)
	return l, nil
}

func createLogFile(logFilePath string) (*os.File, error) {
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, err
	}
	return os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
}

func (l *Logger) logInitialization(config *Config) {
	l.Info("Logger initialized",
		"log_file", config.LogFile,
		"log_level", config.LogLevel,
		"console_output", config.ConsoleOut,
		"file_output", config.FileOut)
}

// toSlogArgs converts key-value pairs to slog-friendly args (same format: key, value, ...).
func toSlogArgs(fields ...any) []any { return fields }

// Debug logs a debug message.
func (l *Logger) Debug(msg string, fields ...any) {
	if !l.shouldLog(msg, fields...) {
		return
	}
	l.slog.Debug(msg, toSlogArgs(fields...)...)
}

// Info logs an info message.
func (l *Logger) Info(msg string, fields ...any) {
	if !l.shouldLog(msg, fields...) {
		return
	}
	l.slog.Info(msg, toSlogArgs(fields...)...)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, fields ...any) {
	if !l.shouldLog(msg, fields...) {
		return
	}
	l.slog.Warn(msg, toSlogArgs(fields...)...)
}

// Error logs an error message. Errors are never deduplicated.
func (l *Logger) Error(msg string, fields ...any) {
	l.slog.Error(msg, toSlogArgs(fields...)...)
}

// DebugContext logs a debug message with context (context is not yet used for correlation).
func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...any) {
	if !l.shouldLog(msg, fields...) {
		return
	}
	l.slog.DebugContext(ctx, msg, toSlogArgs(fields...)...)
}

// InfoContext logs an info message with context.
func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...any) {
	if !l.shouldLog(msg, fields...) {
		return
	}
	l.slog.InfoContext(ctx, msg, toSlogArgs(fields...)...)
}

// WarnContext logs a warning message with context.
func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...any) {
	if !l.shouldLog(msg, fields...) {
		return
	}
	l.slog.WarnContext(ctx, msg, toSlogArgs(fields...)...)
}

// ErrorContext logs an error message with context.
func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...any) {
	l.slog.ErrorContext(ctx, msg, toSlogArgs(fields...)...)
}

// Sync flushes any buffered logs. If file output is used, syncs the file.
func (l *Logger) Sync() error {
	if l.file != nil {
		return l.file.Sync()
	}
	return nil
}
