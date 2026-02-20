package logger

// shouldLog determines if a message should be logged based on deduplication rules.
// Dedup key is msg + first field value (if the first field is a string). Same key and
// equal field values skip logging. Errors are never deduplicated (caller does not use shouldLog).
// lastLogs is unbounded; acceptable for a single process/game session.
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
