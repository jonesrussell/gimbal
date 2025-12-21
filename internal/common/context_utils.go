package common

import "context"

// CheckContextCancellation checks if the context is cancelled and returns an error if so.
// This helper eliminates duplicated cancellation checks across systems.
func CheckContextCancellation(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
