package linkchecker

import (
	"context"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"time"
)

type Result struct {
	URL      string
	Status   int
	Location string
	Duration time.Duration
	Err      error
}

const retryDelay = 100

func CheckLink(ctx context.Context, client *http.Client, url string) Result {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{URL: url, Err: fmt.Errorf("request create error: %w", err)}
	}
	req.Header.Set("User-Agent", "linkchecker/1.0")
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return Result{URL: url, Err: err}
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body) // игнорируем ошибку
	duration := time.Since(start)
	location := ""
	loc := resp.Header.Get("Location")
	if loc != "" {
		if u, err := resp.Request.URL.Parse(loc); err == nil {
			location = u.String()
		} else {
			location = loc
		}
	}
	return Result{URL: url, Status: resp.StatusCode, Location: location, Duration: duration, Err: nil}
}

func CheckWithRetry(ctx context.Context, client *http.Client, url string, attempts int) Result {
	delay := retryDelay
	for i := 0; i < attempts; i++ {
		result := CheckLink(ctx, client, url)
		if result.Err == nil && result.Status < 500 && result.Status != 429 {
			return result
		}

		if i == attempts-1 {
			return result // попытки кончились
		}
		jitter := time.Duration(rand.IntN(delay/5)) * time.Millisecond
		wait := time.Duration(delay)*time.Millisecond + jitter
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return Result{URL: url, Err: ctx.Err()}
		}
		delay = delay * 2
	}
	return Result{} // никогда не достигнет
}
