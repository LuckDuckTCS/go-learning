package linkchecker

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestCheckWithRetry(t *testing.T) {
	tests := []struct {
		name       string
		handler    func(calls *atomic.Int32) http.HandlerFunc
		attempts   int
		wantCalls  int32
		wantStatus int
		wantErr    bool
	}{
		{
			name: "success no retry",
			handler: func(calls *atomic.Int32) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					w.WriteHeader(http.StatusOK)
				}
			},
			attempts:   3,
			wantCalls:  1,
			wantStatus: http.StatusOK,
		},
		{
			name: "404 no retry",
			handler: func(calls *atomic.Int32) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					w.WriteHeader(http.StatusNotFound)
				}
			},
			attempts:   3,
			wantCalls:  1,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "503 exhausts attempts",
			handler: func(calls *atomic.Int32) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					w.WriteHeader(http.StatusServiceUnavailable)
				}
			},
			attempts:   3,
			wantCalls:  3,
			wantStatus: http.StatusServiceUnavailable,
		},
		{
			name: "503 then success",
			handler: func(calls *atomic.Int32) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					if calls.Add(1) <= 2 {
						w.WriteHeader(http.StatusServiceUnavailable)
						return
					}
					w.WriteHeader(http.StatusOK)
				}
			},
			attempts:   3,
			wantCalls:  3,
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32

			srv := httptest.NewServer(tc.handler(&calls))
			defer srv.Close()

			client := &http.Client{Timeout: 2 * time.Second}

			got := CheckWithRetry(context.Background(), client, srv.URL, tc.attempts)

			if calls.Load() != tc.wantCalls {
				t.Errorf("calls = %d, want %d", calls.Load(), tc.wantCalls)
			}
			if got.Status != tc.wantStatus {
				t.Errorf("Status = %d, want %d", got.Status, tc.wantStatus)
			}
			if (got.Err != nil) != tc.wantErr {
				t.Errorf("Err = %v, wantErr %v", got.Err, tc.wantErr)
			}
		})
	}
}

func TestCheckLinkTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 50 * time.Millisecond}
	res := CheckLink(context.Background(), client, srv.URL)

	// res.Err != nil
	// проверить, что это именно таймаут
	if !errors.Is(res.Err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want DeadlineExceeded", res.Err)
	}
}

func TestCheckWithRetryCancel(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable) // чтобы пошли ретраи
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	client := &http.Client{Timeout: 5 * time.Second}
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	res := CheckWithRetry(ctx, client, srv.URL, 10)
	elapsed := time.Since(start)

	// проверить: res.Err — context.Canceled
	if !errors.Is(res.Err, context.Canceled) {
		t.Fatalf("err = %v, want Canceled", res.Err)
	}
	// проверить: elapsed заметно меньше, чем 10 попыток с задержками
	if elapsed > time.Second {
		t.Fatalf("elapsed = %v, want < 1s (cancellation should stop retries)", elapsed)
	}
	// проверить: calls меньше 10
	if calls.Load() >= 10 {
		t.Fatalf("calls = %d, want < 10 (retries should stop after cancel)", calls.Load())
	}
}
