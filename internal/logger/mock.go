package logger

// MockLogger is a no-op logger for testing
type MockLogger struct{}

// NewMock creates a new mock logger
func NewMock() *MockLogger {
	return &MockLogger{}
}

// Debug implements Logger interface
func (m *MockLogger) Debug(msg string, fields ...any) {}

// Info implements Logger interface
func (m *MockLogger) Info(msg string, fields ...any) {}

// Warn implements Logger interface
func (m *MockLogger) Warn(msg string, fields ...any) {}

// Error implements Logger interface
func (m *MockLogger) Error(msg string, fields ...any) {}

// Sync implements Logger interface
func (m *MockLogger) Sync() error { return nil }
