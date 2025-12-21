package logger

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

