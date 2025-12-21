package common

import (
	"context"
	"testing"
	"time"
)

func TestCheckContextCancellation(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() context.Context
		wantErr bool
	}{
		{
			name: "active context",
			setup: func() context.Context {
				return context.Background()
			},
			wantErr: false,
		},
		{
			name: "canceled context",
			setup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
		{
			name: "context with timeout (not expired)",
			setup: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
				return ctx
			},
			wantErr: false,
		},
		{
			name: "context with timeout (expired)",
			setup: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(1 * time.Millisecond) // Ensure timeout has occurred
				cancel()
				return ctx
			},
			wantErr: true,
		},
		{
			name: "context with deadline (not expired)",
			setup: func() context.Context {
				ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
				return ctx
			},
			wantErr: false,
		},
		{
			name: "context with deadline (expired)",
			setup: func() context.Context {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
				cancel()
				return ctx
			},
			wantErr: true,
		},
		{
			name: "context with value (still active)",
			setup: func() context.Context {
				return context.WithValue(context.Background(), "key", "value")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			err := CheckContextCancellation(ctx)

			if tt.wantErr && err == nil {
				t.Error("CheckContextCancellation() expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("CheckContextCancellation() unexpected error: %v", err)
			}
			if tt.wantErr && err != nil && ctx.Err() != err {
				t.Errorf("CheckContextCancellation() error = %v, want %v", err, ctx.Err())
			}
		})
	}
}

func TestCheckContextCancellation_ErrorType(t *testing.T) {
	// Test that canceled context returns context.Canceled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := CheckContextCancellation(ctx)
	if err != context.Canceled {
		t.Errorf("CheckContextCancellation() error = %v, want context.Canceled", err)
	}

	// Test that timeout context returns context.DeadlineExceeded
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	time.Sleep(1 * time.Millisecond) // Ensure timeout has occurred
	timeoutCancel()
	timeoutErr := CheckContextCancellation(timeoutCtx)
	if timeoutErr != context.DeadlineExceeded {
		t.Errorf("CheckContextCancellation() timeout error = %v, want context.DeadlineExceeded", timeoutErr)
	}
}

func TestCheckContextCancellation_Concurrent(t *testing.T) {
	// Test that the function is safe for concurrent use
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run multiple goroutines checking the same context
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			err := CheckContextCancellation(ctx)
			if err != nil {
				done <- false
				return
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		if !<-done {
			t.Error("CheckContextCancellation() should not return error for active context")
		}
	}
}
